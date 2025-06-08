from nornir import InitNornir
from nornir.core.filter import F
from tasks.config import config_intf
from tasks.get import get_version, get_interfaces, get_running_config
from rich.console import Console
from utils.display_results import custom_print
from utils.commands import execute_all
from nornir_utils.plugins.functions import print_result
from utils.database.db_conn import db

console = Console()

def main():
    db.init_db()
    nr = InitNornir(config_file="config/config.yaml")
    leaf_switches = nr.filter(F(groups__contains="leaf"))
    
    ### multiple commands ###
    # commands = [
    #     'show version',
    #     'show interfaces',
    #     'show running-config'
    # ]
    # multi_command_results = leaf_switches.run(
    #     task=execute_all,
    #     commands=commands
    # )
    # custom_print(multi_command_results)
    # # print_result(multi_command_results) # print result from nornir_utils

    ### single command ###
    print("\n=== show version ===")
    command_result = leaf_switches.run(task=get_version, save_result=True)
    custom_print(command_result)
    print("\n=== show interfaces ===")
    command_result = leaf_switches.run(task=get_interfaces, save_result=True)
    custom_print(command_result)
    
    ### interface config example ###
    # interface_config = {
    #     "Ethernet1": {
    #         "description": "test",
    #         "no shutdown": "",
    #     }
    # }
    
    # print("\n=== interface config change ===")
    # config_change_results = leaf_switches.run(
    #     task=config_intf,
    #     interfaces=interface_config
    # )
    # custom_print(config_change_results)

if __name__ == "__main__":
    main()
