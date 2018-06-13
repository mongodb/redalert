"""Contains the CheckFailure exception raised by checks which find incorrect system state."""

class CheckFailure(Exception):
    """Raised by any check which finds an incorrect system state."""
    pass
