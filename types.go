package rewire

// Request is a wrapper for a generic request payload
type Request struct {
	Type      ReqType `struc:"int8"`
	Encrypted bool    `struc:"bool"`
	Size      int     `struc:"int32,sizeof=Payload"`
	Payload   []byte  `struc:"[]int16"`
}
