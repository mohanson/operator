use actix_web::{get, post, web, App, HttpServer, Responder};

mod datatype;
use datatype::*;

#[get("/hello")]
async fn hello() -> impl Responder {
    let r = serde_json::to_string(&HostInfo {
        pubkey: include_str!("../res/alice.pub").trim_end().to_string(),
        ckb_address: "ckb1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsqwau7qpcpealv6xf3a37pdcq6ajhwuyaxgs5g955".into(),
    });
    r.unwrap()
}

#[post("/create_room")]
async fn create_room(body: web::Json<RoomInfo>) -> impl Responder {
    println!("{:?}", body);
    "Ok"
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("listen and server on 127.0.0.1:8080");
    HttpServer::new(|| App::new().service(hello).service(create_room))
        .bind(("127.0.0.1", 8080))?
        .run()
        .await
}
