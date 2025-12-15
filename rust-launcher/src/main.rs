#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use dpi::{LogicalPosition, LogicalSize, PhysicalSize};

use dotenvy::from_filename;
use std::fs::File;
use std::path::PathBuf;
use std::process::{Child, Command};

use std::process::Stdio;
use std::thread;
use std::time::Duration;
use winit::{
    application::ApplicationHandler,
    event::WindowEvent,
    event_loop::{ActiveEventLoop, EventLoop},
    platform::windows::IconExtWindows,
    window::{Icon, Window, WindowId},
};
use wry::{Rect, WebViewBuilder};
const DEFAULT_USER_DATA_FOLDER: &str = "data/webview2";
const RESOURCES_FOLDER: &str = "data/resources";

const DEV_MODE: bool = cfg!(debug_assertions);

const DEV_ENV_PATH: &str = "./../go-backend/.env";

const PROD_ENV_PATH: &str = ".env";

#[derive(Default)]
struct State {
    window: Option<Window>,
    webview: Option<wry::WebView>,
}

fn get_window_icon(desired_size: u32) -> Option<Icon> {
    let icon_path = get_icon_path();
    if icon_path.exists() {
        match Icon::from_path(
            icon_path,
            Some(PhysicalSize::new(desired_size, desired_size)),
        ) {
            Ok(icon) => return Some(icon),
            Err(e) => {
                eprintln!("⚠️ Не удалось загрузить иконку из файла: {}", e);
            }
        }
    }
    None
}

fn get_icon_path() -> PathBuf {
    // Определяем путь в зависимости от структуры проекта
    let project_root = if DEV_MODE {
        // В dev-режиме иконка в корне проекта
        PathBuf::from("./../").join("assets")
    } else {
        // В релизе - относительно исполняемого файла
        std::env::current_exe()
            .unwrap()
            .parent()
            .unwrap()
            .join(RESOURCES_FOLDER)
    };

    project_root.join("icon.ico")
}

impl ApplicationHandler for State {
    fn resumed(&mut self, event_loop: &ActiveEventLoop) {
        unsafe {
            std::env::set_var(
                "WEBVIEW2_USER_DATA_FOLDER",
                std::env::current_exe()
                    .unwrap()
                    .parent()
                    .unwrap()
                    .join(DEFAULT_USER_DATA_FOLDER),
            );
        }
        let mut attributes = Window::default_attributes();
        attributes.visible = false; // что-бы не было черного квадрата при инициализации
        attributes.title = "EVE Traders".to_string();

        attributes.inner_size = Some(LogicalSize::new(1200, 800).into());
        // attributes.maximized = true;
        attributes.decorations = true; //Верхняя панель
        // dbg!(&attributes);
        let window = event_loop.create_window(attributes).unwrap();

        let webview = WebViewBuilder::new()
            .with_devtools(false)
            .with_ipc_handler(move |msg| {
                // Process message
                let message = msg.body();

                match message.as_str() {
                    "-update" => {
                        std::thread::spawn(move || {
                            Command::new(if DEV_MODE {
                                "./../build/updater.exe"
                            } else {
                                "./updater.exe"
                            })
                            .arg("-command")
                            .arg("download")
                            .spawn()
                            .expect("Failed to start update process");
                        });
                    }
                    _ => {
                        println!("no command")
                    }
                };
            })
            .with_initialization_script(
                "
                // async function main() {
                //     for (let i = 0; i < 10; i++) {
                //         window.ipc.postMessage(`запрос из клиента, итерация ${i}`);
                //         await new Promise(resolve => setTimeout(resolve, 1000));
                //     }
                //     alert('закончено');
                // }
                // main();
            ",
            )
            .with_url(
                std::env::var("WEB_SERVER_ADDR").unwrap_or("http://localhost:8080".to_string()),
            )
            .build(&window)
            .unwrap();

        //Окончательная инициализация
        {
            window.set_window_icon(get_window_icon(32));
            //window.set_maximized(true);
            window.set_visible(true);

            self.window = Some(window);
            self.webview = Some(webview);
        }
    }

    fn window_event(
        &mut self,
        _event_loop: &ActiveEventLoop,
        _window_id: WindowId,
        event: WindowEvent,
    ) {
        match event {
            WindowEvent::Resized(size) => {
                let window = self.window.as_ref().unwrap();
                let webview = self.webview.as_ref().unwrap();

                let size = size.to_logical::<u32>(window.scale_factor());
                webview
                    .set_bounds(Rect {
                        position: LogicalPosition::new(0, 0).into(),
                        size: LogicalSize::new(size.width, size.height).into(),
                    })
                    .unwrap();
            }
            WindowEvent::CloseRequested => {
                self.webview = None;
                _event_loop.exit();

                //std::process::exit(0);
            }
            WindowEvent::CursorMoved { position, .. } => {
                println!("Cursor moved to {:?}", position);
            }
            WindowEvent::CursorEntered { .. } => {
                println!("Cursor entered window");
            }

            _ => {}
        }
    }

    fn about_to_wait(&mut self, _event_loop: &ActiveEventLoop) {
        #[cfg(any(
            target_os = "linux",
            target_os = "dragonfly",
            target_os = "freebsd",
            target_os = "netbsd",
            target_os = "openbsd",
        ))]
        {
            while gtk::events_pending() {
                gtk::main_iteration_do(false);
            }
        }
    }
}

struct BackendManager {
    process: Option<Child>,
}

impl BackendManager {
    fn new() -> Self {
        Self { process: None }
    }

    fn start(&mut self) -> bool {
        const MAX_ATTEMPTS: u32 = 10;

        // if !std::path::Path::new("./go-backend.exe").exists() {
        //     eprintln!("Backend executable not found!");
        //     return false;
        // }

        for attempt in 1..=MAX_ATTEMPTS {
            println!("Attempt {} to start backend...", attempt);

            match Command::new("./go-backend.exe")
                .arg("-debug")
                .stdout(Stdio::piped())
                .stderr(Stdio::piped())
                .spawn()
            {
                Ok(mut child) => {
                    // Даем процессу время на запуск
                    thread::sleep(Duration::from_millis(500));

                    // Проверяем, не завершился ли процесс
                    match child.try_wait() {
                        Ok(Some(status)) => {
                            eprintln!("Backend exited with status: {:?}", status);
                            if attempt < MAX_ATTEMPTS {
                                thread::sleep(Duration::from_millis(1000));
                                continue;
                            }
                            return false;
                        }
                        Ok(None) => {
                            println!("Backend started successfully!");
                            self.process = Some(child);
                            return true;
                        }
                        Err(e) => {
                            eprintln!("Error: {}", e);
                        }
                    }
                }
                Err(e) => {
                    eprintln!("Failed to start backend: {}", e);
                    if attempt < MAX_ATTEMPTS {
                        thread::sleep(Duration::from_millis(2000));
                    }
                }
            }
        }

        eprintln!("Failed to start backend after {} attempts", MAX_ATTEMPTS);
        false
    }
}

impl Drop for BackendManager {
    fn drop(&mut self) {
        if let Some(mut child) = self.process.take() {
            let _ = child.kill();
            let _ = child.wait();
            println!("Backend process terminated.");
        }
    }
}

fn main() -> wry::Result<()> {
    #[cfg(any(
        target_os = "linux",
        target_os = "dragonfly",
        target_os = "freebsd",
        target_os = "netbsd",
        target_os = "openbsd",
    ))]
    {
        use gtk::prelude::DisplayExtManual;

        gtk::init().unwrap();
        if gtk::gdk::Display::default().unwrap().backend().is_wayland() {
            panic!("This example doesn't support wayland!");
        }

        winit::platform::x11::register_xlib_error_hook(Box::new(|_display, error| {
            let error = error as *mut x11_dl::xlib::XErrorEvent;
            (unsafe { (*error).error_code }) == 170
        }));
    }

    from_filename(if DEV_MODE {
        DEV_ENV_PATH
    } else {
        PROD_ENV_PATH
    })
    .expect("Failed to load env file");

    // Updater
    if DEV_MODE {
        unsafe {
            std::env::set_var("WEB_SERVER_ADDR", "http://localhost:5173");
        }
    } else {
        if File::open("data/config.json").is_err() {
            Command::new(if DEV_MODE {
                "./../build/updater.exe"
            } else {
                "./updater.exe"
            })
            .arg("-command")
            .arg("download")
            .spawn()
            .expect("Failed to start update process");
        }
    }

    let event_loop = EventLoop::new().unwrap();

    let mut state = State::default();

    let mut backend_manager = BackendManager::new();
    if !backend_manager.start() {
        eprintln!("Failed to start backend. Exiting.");
        return Ok(());
    }
    event_loop.run_app(&mut state).unwrap();

    Ok(())
}
