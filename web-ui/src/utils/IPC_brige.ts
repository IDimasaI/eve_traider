declare global {
    interface Window {
        ipc: {
            postMessage: (message: string) => void
        }
    }
}

export const ipc_update_app = () => window.ipc.postMessage(`-update`)

