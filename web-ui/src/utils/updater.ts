import { ipc_update_app } from "./IPC_brige"

type Status = {
    Current_Version: string;
    Status: string;
    Progress: string;
    Timestamp: number;
}

export class Updater {
    public current_version: string
    private pulling_delay = 200
    constructor() {
        this.current_version = ""
    }

    public async update_app(): Promise<void> {

        if (await this.check_update() == false) {
            return
        }
        console.log("Новая версия")
        let status = await this.get_status()
        ipc_update_app()
        while (
            status["Status"] !== "finished" &&
            status["Status"] !== "error") {

            console.log(status["Status"]);
            await new Promise(resolve => setTimeout(resolve, this.pulling_delay));
            status = await this.get_status();
        }

        if (status["Status"] == "finished") {
            console.log("Обновление завершено")
            this.current_version = status["Current_Version"]
            window.location.reload()
        }
        return
    }

    private async get_status(): Promise<Status> {
        return (await (await fetch("/api/v2/update_status")).json())
    }

    public async check_update(): Promise<boolean> {

        if (this.current_version == "" || this.current_version == null) {
            await this.Init_current_version()
        }

        if (this.current_version != await this.get_latest_gh_version()) {
            return true
        }


        console.log("Новая версия не обнаружена")
        return false
    }

    private async Init_current_version(): Promise<void> {
        this.current_version = await this.get_current_version()
    }

    // private set_current_version(version: string): void {
    //     this.current_version = version
    // }

    private async get_current_version(): Promise<string> {
        return (await (await fetch("/api/v2/update_status")).json())["Current_Version"]
    }

    private async get_latest_gh_version(): Promise<string> {
        return (await (await fetch("https://api.github.com/repos/IDimasaI/eve_traider/releases/latest")).json())["tag_name"]
    }

}