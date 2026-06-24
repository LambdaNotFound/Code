import unittest

from interview.rippling_file_system import FileSystem


class TestLevel1(unittest.TestCase):
    def setUp(self) -> None:
        self.fs = FileSystem()

    def test_create_file_at_root(self) -> None:
        self.assertTrue(self.fs.create_file("/", "readme.txt"))

    def test_create_directory_at_root(self) -> None:
        self.assertTrue(self.fs.create_directory("/", "docs"))

    def test_list_contents_root_empty(self) -> None:
        self.assertEqual(self.fs.list_contents("/"), [])

    def test_list_contents_root_with_entries(self) -> None:
        self.fs.create_file("/", "b.txt")
        self.fs.create_directory("/", "a")
        self.assertEqual(self.fs.list_contents("/"), ["a", "b.txt"])

    def test_delete_file_from_root(self) -> None:
        self.fs.create_file("/", "readme.txt")
        self.assertTrue(self.fs.delete("/readme.txt"))
        self.assertEqual(self.fs.list_contents("/"), [])

    def test_delete_directory_from_root(self) -> None:
        self.fs.create_directory("/", "docs")
        self.assertTrue(self.fs.delete("/docs"))
        self.assertEqual(self.fs.list_contents("/"), [])


class TestLevel2(unittest.TestCase):
    def setUp(self) -> None:
        self.fs = FileSystem()
        self.fs.create_directory("/", "src")

    def test_create_file_in_level1_dir(self) -> None:
        self.assertTrue(self.fs.create_file("/src", "main.py"))

    def test_create_directory_in_level1_dir(self) -> None:
        self.assertTrue(self.fs.create_directory("/src", "utils"))

    def test_list_contents_level1_dir(self) -> None:
        self.fs.create_file("/src", "main.py")
        self.fs.create_directory("/src", "utils")
        self.assertEqual(self.fs.list_contents("/src"), ["main.py", "utils"])

    def test_delete_file_from_level1_dir(self) -> None:
        self.fs.create_file("/src", "main.py")
        self.assertTrue(self.fs.delete("/src/main.py"))
        self.assertEqual(self.fs.list_contents("/src"), [])

    def test_delete_empty_dir_from_level1_dir(self) -> None:
        self.fs.create_directory("/src", "utils")
        self.assertTrue(self.fs.delete("/src/utils"))
        self.assertEqual(self.fs.list_contents("/src"), [])


class TestEdgeCases(unittest.TestCase):
    def setUp(self) -> None:
        self.fs = FileSystem()

    def test_delete_nonexistent_path_returns_false(self) -> None:
        self.assertFalse(self.fs.delete("/nonexistent"))

    def test_delete_nonempty_directory_returns_false(self) -> None:
        self.fs.create_directory("/", "docs")
        self.fs.create_file("/docs", "readme.txt")
        self.assertFalse(self.fs.delete("/docs"))

    def test_delete_root_returns_false(self) -> None:
        self.assertFalse(self.fs.delete("/"))

    def test_create_file_in_nonexistent_path_returns_false(self) -> None:
        self.assertFalse(self.fs.create_file("/nonexistent", "file.txt"))

    def test_create_directory_in_nonexistent_path_returns_false(self) -> None:
        self.assertFalse(self.fs.create_directory("/nonexistent", "dir"))

    def test_create_file_in_file_returns_false(self) -> None:
        self.fs.create_file("/", "readme.txt")
        self.assertFalse(self.fs.create_file("/readme.txt", "nested.txt"))

    def test_duplicate_file_create_returns_false(self) -> None:
        self.fs.create_file("/", "readme.txt")
        self.assertFalse(self.fs.create_file("/", "readme.txt"))

    def test_duplicate_directory_create_returns_false(self) -> None:
        self.fs.create_directory("/", "docs")
        self.assertFalse(self.fs.create_directory("/", "docs"))

    def test_list_contents_on_file_returns_empty(self) -> None:
        self.fs.create_file("/", "readme.txt")
        self.assertEqual(self.fs.list_contents("/readme.txt"), [])

    def test_list_contents_sorted_alphabetically(self) -> None:
        self.fs.create_file("/", "z.txt")
        self.fs.create_file("/", "a.txt")
        self.fs.create_directory("/", "m")
        self.assertEqual(self.fs.list_contents("/"), ["a.txt", "m", "z.txt"])


if __name__ == "__main__":
    unittest.main()
