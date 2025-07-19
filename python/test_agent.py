import pytest
from agent import run_agent

def test_run_agent_summarize():
    prompt = "summarize this document"
    expected_response = "This page is a sample domain used for illustrative purposes."
    assert run_agent(prompt) == expected_response

def test_run_agent_general_prompt():
    prompt = "hello world"
    expected_response = "[Agent Reply] I received: hello world"
    assert run_agent(prompt) == expected_response

def test_run_agent_empty_prompt():
    prompt = ""
    expected_response = "[Agent Reply] I received: "
    assert run_agent(prompt) == expected_response
