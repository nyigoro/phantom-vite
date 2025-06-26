import sys

def run_agent(prompt):
    print(f"[Phantom Agent] Prompt received: {prompt}")
    # Placeholder response - real AI model integration comes next
    if "summarize" in prompt:
        return "This page is a sample domain used for illustrative purposes."
    return f"[Agent Reply] I received: {prompt}"

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python agent.py <prompt>")
        sys.exit(1)
    
    prompt = " ".join(sys.argv[1:])
    response = run_agent(prompt)
    print(response)
