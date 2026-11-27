//! LeetCode's node types, shared by every solution under `dsa/`; mirrors
//! `go/dsa/types/`. The definitions are LeetCode's own Rust ones.

mod list;
mod tree;

pub use list::ListNode;
pub use tree::TreeNode;
