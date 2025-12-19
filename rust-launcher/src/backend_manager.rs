use std::process::{Child, Command};

use std::process::Stdio;
use std::thread;
use std::time::Duration;
pub struct BackendManager {
    process: Option<Child>,
}

impl BackendManager {
    pub fn new() -> Self {
        Self { process: None }
    }

    pub fn start(&mut self) -> bool {
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
