import logging
from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Optional

logger = logging.getLogger(__name__)

# Design:
#   _consumer_service  ConsumerService  — userId -> Consumer; hard dependency, failure aborts the whole profile
#   _payment_service   PaymentService   — consumerId -> PaymentInfo; soft dependency, failure leaves fields None
#   _address_service   AddressService   — consumerId -> Address; soft dependency, failure leaves field None
#
# Key decision: each collaborator is injected as an ABC rather than a concrete class, so tests can
# swap in fakes that raise on demand. Consumer lookup is a hard dependency (no consumerId, no profile
# at all) while Payment/Address are wrapped in their own try/except (bulkhead isolation) so one
# downstream failure never blanks out a field the other call already fetched successfully.


@dataclass
class Consumer:
    id: str
    name: str


@dataclass
class PaymentInfo:
    default_method: str
    gift_card_balance: float


@dataclass
class Address:
    line1: str
    city: str
    zip: str


@dataclass
class UserProfile:
    consumer_id: str
    name: str
    default_payment_method: Optional[str] = None
    gift_card_balance: Optional[float] = None
    address: Optional[Address] = None


class UserNotFoundError(Exception):
    pass


class ConsumerService(ABC):
    @abstractmethod
    def get_consumer(self, user_id: str) -> Consumer: ...


class PaymentService(ABC):
    @abstractmethod
    def get_payment_info(self, consumer_id: str) -> PaymentInfo: ...


class AddressService(ABC):
    @abstractmethod
    def get_address(self, consumer_id: str) -> Address: ...


class BootstrapService:
    def __init__(
        self,
        consumer_service: ConsumerService,
        payment_service: PaymentService,
        address_service: AddressService,
    ):
        self._consumer_service = consumer_service
        self._payment_service = payment_service
        self._address_service = address_service

    def get_user_profile(self, user_id: str) -> UserProfile:  # T: O(1) — three downstream calls, S: O(1)
        try:
            consumer = self._consumer_service.get_consumer(user_id)
        except Exception as e:
            logger.error("ConsumerService failed for userId=%s: %s", user_id, e)
            raise UserNotFoundError(f"Could not resolve user {user_id}") from e

        profile = UserProfile(consumer_id=consumer.id, name=consumer.name)

        try:
            payment_info = self._payment_service.get_payment_info(consumer.id)
            profile.default_payment_method = payment_info.default_method
            profile.gift_card_balance = payment_info.gift_card_balance
        except Exception as e:
            logger.error("PaymentService failed for consumerId=%s: %s", consumer.id, e)

        try:
            profile.address = self._address_service.get_address(consumer.id)
        except Exception as e:
            logger.error("AddressService failed for consumerId=%s: %s", consumer.id, e)

        return profile
