import sys, os
sys.path.insert(0, os.path.dirname(__file__))

import pytest
from doordash_bootstrap import (
    BootstrapService,
    ConsumerService,
    PaymentService,
    AddressService,
    Consumer,
    PaymentInfo,
    Address,
    UserProfile,
    UserNotFoundError,
)


class _FakeConsumerService(ConsumerService):
    def __init__(self, consumer: Consumer = None, error: Exception = None):
        self._consumer = consumer
        self._error = error

    def get_consumer(self, user_id: str) -> Consumer:
        if self._error:
            raise self._error
        return self._consumer


class _FakePaymentService(PaymentService):
    def __init__(self, payment_info: PaymentInfo = None, error: Exception = None):
        self._payment_info = payment_info
        self._error = error

    def get_payment_info(self, consumer_id: str) -> PaymentInfo:
        if self._error:
            raise self._error
        return self._payment_info


class _FakeAddressService(AddressService):
    def __init__(self, address: Address = None, error: Exception = None):
        self._address = address
        self._error = error

    def get_address(self, consumer_id: str) -> Address:
        if self._error:
            raise self._error
        return self._address


@pytest.fixture
def consumer():
    return Consumer(id="123", name="Alice")


@pytest.fixture
def payment_info():
    return PaymentInfo(default_method="Credit Card", gift_card_balance=50.0)


@pytest.fixture
def address():
    return Address(line1="123 Main St", city="Anytown", zip="12345")


class TestHappyPath:
    def test_assembles_full_profile(self, consumer, payment_info, address):
        service = BootstrapService(
            _FakeConsumerService(consumer=consumer),
            _FakePaymentService(payment_info=payment_info),
            _FakeAddressService(address=address),
        )

        profile = service.get_user_profile("user123")

        assert profile == UserProfile(
            consumer_id="123",
            name="Alice",
            default_payment_method="Credit Card",
            gift_card_balance=50.0,
            address=Address(line1="123 Main St", city="Anytown", zip="12345"),
        )


class TestConsumerServiceFailure:
    def test_invalid_user_raises(self, payment_info, address):
        service = BootstrapService(
            _FakeConsumerService(error=RuntimeError("not found")),
            _FakePaymentService(payment_info=payment_info),
            _FakeAddressService(address=address),
        )

        with pytest.raises(UserNotFoundError):
            service.get_user_profile("bad_user")


class TestPaymentServiceFailure:
    def test_partial_profile_without_payment_fields(self, consumer, address):
        service = BootstrapService(
            _FakeConsumerService(consumer=consumer),
            _FakePaymentService(error=RuntimeError("payment down")),
            _FakeAddressService(address=address),
        )

        profile = service.get_user_profile("user123")

        assert profile.consumer_id == "123"
        assert profile.name == "Alice"
        assert profile.default_payment_method is None
        assert profile.gift_card_balance is None
        assert profile.address == address


class TestAddressServiceFailure:
    def test_partial_profile_without_address(self, consumer, payment_info):
        service = BootstrapService(
            _FakeConsumerService(consumer=consumer),
            _FakePaymentService(payment_info=payment_info),
            _FakeAddressService(error=RuntimeError("address down")),
        )

        profile = service.get_user_profile("user123")

        assert profile.consumer_id == "123"
        assert profile.name == "Alice"
        assert profile.default_payment_method == "Credit Card"
        assert profile.gift_card_balance == 50.0
        assert profile.address is None


class TestMultipleServiceFailures:
    def test_only_consumer_fields_present(self, consumer):
        service = BootstrapService(
            _FakeConsumerService(consumer=consumer),
            _FakePaymentService(error=RuntimeError("payment down")),
            _FakeAddressService(error=RuntimeError("address down")),
        )

        profile = service.get_user_profile("user123")

        assert profile == UserProfile(consumer_id="123", name="Alice")
