use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize)]
pub struct HostInfo {
    pub ckb_address: String,
    pub pubkey: String,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct RoomInfo {
    pub host_ckb_address: String,
    pub room_owner_pubkey: String,
    pub room_member_pubkey_list: Vec<String>,
    pub timelock: u64,
    pub transfer_price: u64,
    pub data_keep_price: u64,
    pub current_charge_points: u64,
}
