from __future__ import annotations
from typing import Optional


class Node:
    def __init__(self, name: str, is_directory: bool, parent: Optional[Node] = None):
        self.name = name
        self.is_directory = is_directory
        self.children: dict[str, Node] = {}
        self.parent: Optional[Node] = parent


class FileSystem:
    def __init__(self) -> None:
        self._root = Node("/", is_directory=True)

    def _traverse(self, path: str) -> Optional[Node]:
        segments = [s for s in path.split("/") if s]
        node = self._root
        for segment in segments:
            if segment not in node.children:
                return None
            node = node.children[segment]
        return node

    def create_file(self, path: str, name: str) -> bool:
        parent = self._traverse(path)
        if parent is None or not parent.is_directory or name in parent.children:
            return False
        parent.children[name] = Node(name, is_directory=False, parent=parent)
        return True

    def create_directory(self, path: str, name: str) -> bool:
        parent = self._traverse(path)
        if parent is None or not parent.is_directory or name in parent.children:
            return False
        parent.children[name] = Node(name, is_directory=True, parent=parent)
        return True

    def delete(self, path: str) -> bool:
        node = self._traverse(path)
        if node is None or node.parent is None or node.children:
            return False
        del node.parent.children[node.name]
        return True

    def list_contents(self, path: str) -> list[str]:
        node = self._traverse(path)
        if node is None or not node.is_directory:
            return []
        return sorted(node.children.keys())
