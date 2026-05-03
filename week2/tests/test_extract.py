import os
import pytest

from unittest.mock import patch

from ..app.services import extract
from ..app.services.extract import extract_action_items, extract_action_items_llm


def test_extract_action_items_with_bullets_and_checkboxes_returns_correct_items():
    text = """
    Notes from meeting:
    - [ ] Set up database
    * implement API extract endpoint
    1. Write tests
    Some narrative sentence.
    """.strip()

    items = extract_action_items(text)
    assert "Set up database" in items
    assert "implement API extract endpoint" in items
    assert "Write tests" in items


@patch.object(extract, "chat")
def test_extract_action_items_llm_with_bullets_and_checkboxes_returns_correct_items(mock_chat):
    mock_chat.return_value = {
        "message": {
            "content": '{"action_items": ["Set up database", "implement API extract endpoint", "Write tests"]}'
        }
    }
    
    text = """
    Notes from meeting:
    - [ ] Set up database
    * implement API extract endpoint
    1. Write tests
    Some narrative sentence.
    """.strip()

    items = extract_action_items_llm(text)
    assert "Set up database" in items
    assert "implement API extract endpoint" in items
    assert "Write tests" in items
    
    # Verify the LLM was called with the input text
    called_args, called_kwargs = mock_chat.call_args
    assert "messages" in called_kwargs
    assert text in called_kwargs["messages"][0]["content"]


@patch.object(extract, "chat")
def test_extract_action_items_llm_with_keyword_prefixed_returns_correct_items(mock_chat):
    mock_chat.return_value = {
        "message": {
            "content": '{"action_items": ["Review PR", "Email client"]}'
        }
    }
    text = "TODO: Review PR\nACTION: Email client"
    items = extract_action_items_llm(text)
    assert items == ["Review PR", "Email client"]
    
    called_args, called_kwargs = mock_chat.call_args
    assert text in called_kwargs["messages"][0]["content"]


@patch.object(extract, "chat")
def test_extract_action_items_llm_with_empty_input_returns_empty_list(mock_chat):
    mock_chat.return_value = {
        "message": {
            "content": '{"action_items": []}'
        }
    }
    items = extract_action_items_llm("")
    assert items == []
