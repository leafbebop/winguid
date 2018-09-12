package guid

import (
	"testing"
	"fmt"
)

func TestDehex(t *testing.T) {
	v,err:= dehex(byte('c'),0)
	if v!=12 || err!=nil {
		fmt.Println("dehex 12",v,err)
		t.Fatal("c")
	}

	v,err=dehex(byte('F'),0)

	if v!=15 || err!=nil {
		fmt.Println("dehex15",v,err)
		t.Fatal("F")
	}

	v,err=dehex(byte('7'),0)

	if v!=7 || err!=nil {
		fmt.Println("dehex7",v,err)
		t.Fatal("7")
	}


	v,err=dehex(byte('-'),17)

	if err == nil {
		fmt.Println("dehex -")
		t.Fatal("-")
	}

	e,ok:=err.(ErrBadFormat)
	e0:=ErrBadFormat{'-',17}
	if !ok || e != e0 {
		fmt.Println("dehex -",v,err)
		t.Fatal("-,17")
	}
}

type stringCases struct {
	g GUID
	s string
}

var  sC = []stringCases {
	{ GUID{ 1,2,3,[8]uint8{1,1,1,1,1,1,1,1} },
		"{00000001-0002-0003-0101-010101010101}",
	}, 
	{ GUID{ 0x10000001,0xb,0x300,[8]uint8{1,2,3,10,21,1,0xaa,0xbb} },
		"{10000001-000b-0300-0102-030a1501aabb}",
	}, 
	{ GUID{ 0x770aae78,0xf26f,0x4dba,[8]uint8{0xa8, 0x29, 0x25, 0x3c, 0x83, 0xd1, 0xb3, 0x87} },
		"{770aae78-f26f-4dba-a829-253c83d1b387}",
	},
	{ GUID{0xcafcb56c,0x6ac3,0x4889,[8]uint8{0xbf, 0x47, 0x9e, 0x23, 0xbb, 0xd2, 0x60, 0xec} },
		"{cafcb56c-6ac3-4889-bf47-9e23bbd260ec}",
	},
}
	
func TestString(t *testing.T) {
	for _,c:=range sC {
		if p:=c.g.String(); p != c.s {
			fmt.Printf("%#v (%s) -> %s\n",c.g,c.s,p)
			fmt.Println([]byte(p),len(p))
			fmt.Println([]byte(c.s),len(c.s))
			t.Fatal("String")
		}
	}
}

func TestFromString(t *testing.T) {
	for _,c:=range sC {
		p,err:=FromString(c.s)
		if err!=nil {
			fmt.Printf("%s(%#v)",c.s,c.g)
			t.Fatal(err)
		}
		if p!=c.g {
			fmt.Printf("%s(%#v) -> %#v\n",c.s,c.g,p)
			t.Fatal("FromString")
		}
	}
}

func TestITOA(t *testing.T) {
	for i:=0;i<50;i++ {
		if got,want:=itoa(i),fmt.Sprint(i); got!=want {
			fmt.Printf("%v != %v\n",got,want)
			t.Fatal("itoa")
		}
	}
}

