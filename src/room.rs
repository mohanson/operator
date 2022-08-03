use std::borrow::BorrowMut;
use std::collections::HashMap;
use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};
use std::sync::{Arc, Mutex, RwLock};

fn hand(id: u32, stream: Arc<RwLock<TcpStream>>, client: &HashMap<u32, Arc<RwLock<TcpStream>>>) {
    let mut reader = stream.read().unwrap().try_clone().unwrap();
    while true {
        let mut buf: Vec<u8> = vec![0; 8];
        reader.read_exact(&mut buf).unwrap();
        let user_id = u32::from_le_bytes(buf[0..4].try_into().unwrap());
        let message_size = u32::from_le_bytes(buf[4..8].try_into().unwrap());
        buf.resize(8 + message_size as usize, 0);
        reader.read_exact(&mut buf[8..]).unwrap();
        let message = String::from_utf8_lossy(&buf[8..]);

        for (k, v) in client.iter() {
            v.write().unwrap().write_all(&buf).unwrap()
        }

        println!("{:?} {:?} {:?} {:?}", id, user_id, message_size, message);
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let listen = String::from("127.0.0.1:51958");
    let mut client = HashMap::<u32, Arc<RwLock<TcpStream>>>::new();
    let mut id = 0;
    println!("listen and server on {}", listen);
    let listener = TcpListener::bind(&listen[..])?;
    for stream in listener.incoming() {
        let stream = stream.unwrap();
        let stream_lock = Arc::new(RwLock::new(stream));
        client.insert(id, stream_lock.clone());
        let client_me = client.clone();
        std::thread::spawn(move || hand(id, stream_lock, &client_me));
        id += 1;
    }

    Ok(())
}
