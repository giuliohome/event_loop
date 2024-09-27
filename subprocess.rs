use std::process::{Command, exit};
use std::env;
use std::io;

// Simulating a CPU-bound task using u128 to avoid overflow
fn cpu_bound() -> u128 {
    (0..10_000_000u128).map(|i| i * i).sum() // Using u128 here
}

// Function to run the CPU-bound task in a subprocess with custom environment variables
fn run_cpu_bound_in_subprocess() -> io::Result<()> {
    // Spawning a new process, passing "cpu" argument to signal the subprocess to run the task
    let output = Command::new(std::env::current_exe()?)
    .arg("cpu")
    .env("CUSTOM_ENV_VAR", "custom_value") // Set environment variable
    .output()?; // Capture the output

    if !output.status.success() {
        eprintln!("Subprocess failed with status: {}", output.status);
        eprintln!("Subprocess error output: {}", String::from_utf8_lossy(&output.stderr));
    }

    // Printing the output from the subprocess
    println!("Subprocess output: {}", String::from_utf8_lossy(&output.stdout));
    Ok(())
}

// Main function that spawns threads and handles the subprocess execution
fn main() -> io::Result<()> {
    // Checking if we are running in subprocess mode (with "cpu" argument)
    if let Some(arg) = std::env::args().nth(1) {
        if arg == "cpu" {
            cpu();
            exit(0); // Exit after running the subprocess task
        }
    }

    // Run the CPU-bound task in a subprocess
    run_cpu_bound_in_subprocess()?;

    Ok(())
}

// The CPU function that runs when the subprocess is spawned
fn cpu() {
    // Fetch and print the custom environment variable
    match env::var("CUSTOM_ENV_VAR") {
        Ok(custom_env_var) => {
            println!("CUSTOM_ENV_VAR: {}", custom_env_var);
        }
        Err(e) => {
            println!("Error retrieving CUSTOM_ENV_VAR: {}", e);
        }
    }

    // Perform the CPU-bound task
    let result = cpu_bound();
    println!("CPU-bound result: {}", result);
}
