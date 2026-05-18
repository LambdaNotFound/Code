"""
DB schema: 
Field | Type
created_at | unix timestamp
dispute_id | int (unique)
loan_id | int (unique)
event_type | string, enum: "opened", "closed_fraud", "closed_not_fraud"
"""
class StatusCalculator:
    # (some code omitted...)

    def _get_dispute_status(self, dispute_id: int) -> DisputeEventType:
        events = self._get_all_events_for_dispute(dispute_id)
        # lambda
        events = sorted(events, key=lambda x: x.created_at)
        return events[-1].event_type

    def _get_loan_status(self, loan_id: int) -> LoanStatus:
        all_dispute_statuses = self._get_all_dispute_statuses_for_loan(loan_id)
        if DisputeEventType.closed_fraud in all_dispute_statuses:
            return LoanStatus.fraudulent
        if DisputeEventType.opened in all_dispute_statuses:
            return LoanStatus.under_investigation
        return LoanStatus.not_fraudulent


    """
    !!! DO NOT modify the code below !!!
    """

    def get_status_per_loan(self) -> dict[int, LoanStatus]:
        all_loan_ids = [entry.loan_id for entry in self.db]

        return {
            loan_id: self._get_loan_status(loan_id)
            for loan_id in all_loan_ids
        }


if __name__ == "__main__":
    db: list[DisputeEvent] = [DisputeEvent.from_dict(x) for x in json.loads(sys.stdin.read())]
    status_calculator = StatusCalculator(db)

    loans_with_status = dict(sorted(status_calculator.get_status_per_loan().items()))

    for (loan_id, status) in loans_with_status.items():
        print(f"{loan_id} {status.value}")

"""
Your task: Design an event consumer StatusManager which processes a stream of DisputeEvent events and supports a function for retrieving the status of a loan with time complexity O(1).

The events are prepopulated by a JSON list of objects (dictionaries) that we pass in as input. To refresh your memory, this is the structure for the dispute events:
"""

class DisputeEvent:
    @classmethod
    def from_dict(cls, data: dict) -> "DisputeEvent":
        return cls(
            created_at=data["created_at"],
            dispute_id=data["dispute_id"],
            loan_id=data["loan_id"],
            event_type=DisputeEventType(data["event_type"]),
        )

from enum import StrEnum
class LoanStatus(StrEnum):
    under_investigation = "under_investigation"
    fraudulent = "fraudulent"
    not_fraudulent = "not_fraudulent"


class StatusManager:
    def __init__(self):
        self.loans: dict[int, dict[str, int]] = defaultdict(lambda : defaultdict(int))
        self.dispute_events_timestamps: dict[int, int] = defaultdict(int)
        self.dispute_statuses: dict[int, DisputeEventType] = {}
    
    def add_event(self, event: DisputeEvent) -> None:
        dispute_id = event.dispute_id
        loan_id = event.loan_id
        current_status = event.event_type
        created_at = event.created_at
        
        if created_at <= self.dispute_events_timestamps[dispute_id]:
            return
        
        if loan_id not in self.loans:
            return
        
        if dispute_id in self.dispute_statuses.keys():
            old_status = self.dispute_statuses[dispute_id]
            self.loans[loan_id][old_status] -= 1
            
        self.loans[loan_id][current_status] += 1
            
        self.dispute_statuses[dispute_id] = current_status
        self.dispute_events_timestamps[dispute_id] = created_at
        

    def get_loan_status(self, loan_id: int) -> LoanStatus | None:
        if self.loans[loan_id][DisputeEventType.closed_fraud] > 0:
            return LoanStatus.fraudulent
        if self.loans[loan_id][DisputeEventType.opened] > 0:
            return LoanStatus.under_investigation
        return LoanStatus.not_fraudulent

if __name__ == "__main__":
    events: list[DisputeEvent] = [DisputeEvent.from_dict(x) for x in json.loads(sys.stdin.read())]

    status_manager = StatusManager()

    for event in events:
        status_manager.add_event(event)

    loan_ids = sorted({e.loan_id for e in events})

    for loan_id in loan_ids:
        print(f"{loan_id} {status_manager.get_loan_status(loan_id).value}")