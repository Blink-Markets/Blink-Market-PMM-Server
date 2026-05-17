package sign

import "encoding/binary"

// EncodeMessage reproduces, byte-for-byte, the concatenation in
// blink_market rfq.move:91-97 (BCS_MESSAGE_CONTRACT, spec §3.1):
//   market_id(32 raw) | side(u8) | price_bps(u64 LE) | size(u64 LE) |
//   seq_number(u64 LE) | expires_at(u64 LE) | pmm(32 raw)
func EncodeMessage(marketID [32]byte, side uint8, priceBps, size, seq, expiresAt uint64, pmm [32]byte) []byte {
	out := make([]byte, 0, 97)
	out = append(out, marketID[:]...)
	out = append(out, side)
	var u [8]byte
	for _, v := range []uint64{priceBps, size, seq, expiresAt} {
		binary.LittleEndian.PutUint64(u[:], v)
		out = append(out, u[:]...)
	}
	out = append(out, pmm[:]...)
	return out
}
