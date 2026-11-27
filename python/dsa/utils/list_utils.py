from dsa.types import ListNode


def build_list(vals: list[int]) -> "ListNode | None":
    dummy = ListNode()
    cur = dummy
    for v in vals:
        cur.next = ListNode(v)
        cur = cur.next
    return dummy.next


def list_to_vals(node: "ListNode | None") -> list[int]:
    result = []
    while node:
        result.append(node.val)
        node = node.next
    return result
