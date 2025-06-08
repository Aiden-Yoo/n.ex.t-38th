from nornir.core.task import AggregatedResult
from rich.table import Table
from rich.console import Console
import json

console = Console()

def custom_print(results: AggregatedResult):
    table = Table(title="Task Results")
    table.add_column("Host", style="cyan")
    table.add_column("Status", style="green")
    table.add_column("Results", style="yellow")

    for host, result in results.items():
        status = "Success" if not result.failed else "Failed"
        if result.failed:
            output = f"Error: {str(result.exception)}"
        else:
            output = json.dumps(result.result, indent=2, ensure_ascii=False)
        table.add_row(host, status, output)

    console.print(table)
