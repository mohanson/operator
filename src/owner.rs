mod datatype;
use datatype::*;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let http_client = reqwest::blocking::Client::new();

    let host_info_str = http_client
        .get("http://127.0.0.1:8080/hello")
        .send()?
        .text()?;
    let host_info: HostInfo = serde_json::from_str(&host_info_str)?;
    let bob_pubkey = include_str!("../res/bob.pub").trim_end();
    let claire_pubkey = include_str!("../res/claire.pub").trim_end();

    let body = RoomInfo {
        host_ckb_address: host_info.ckb_address,
        room_owner_pubkey: bob_pubkey.into(),
        room_member_pubkey_list: vec![claire_pubkey.into()],
        timelock: 1000,
        transfer_price: 1,
        data_keep_price: 1,
        current_charge_points: 0,
    };
    let r = http_client
        .post("http://127.0.0.1:8080/create_room")
        .body(serde_json::to_string(&body)?)
        .header("content-type", "application/json")
        .send()?
        .text()?;
    println!("{:?}", r);
    Ok(())
}
