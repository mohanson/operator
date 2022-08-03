use std::borrow::BorrowMut;
use std::collections::HashMap;
use std::net::{TcpListener, TcpStream};
use std::sync::{Arc, Mutex};

fn hand(id: u32, stream: &TcpStream, client: &HashMap<u32, Arc<TcpStream>>) {
    let mut reader = stream.try_clone().unwrap();
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let listen = String::from("127.0.0.1:51958");
    let mut client = HashMap::<u32, Arc<TcpStream>>::new();
    let mut id = 0;
    println!("listen and server on {}", listen);
    let listener = TcpListener::bind(&listen[..])?;
    for stream in listener.incoming() {
        let stream = stream.unwrap();
        let stream = Arc::new(stream);
        client.insert(id, stream.clone());
        let client_me = client.clone();
        std::thread::spawn(move || hand(id, &stream, &client_me));
        id += 1;
    }

    Ok(())
}
