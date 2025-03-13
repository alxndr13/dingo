// Code generated by cuelang.org/go/pkg/gen. DO NOT EDIT.

package hex

import (
	"cuelang.org/go/internal/core/adt"
	"cuelang.org/go/internal/pkg"
)

func init() {
	pkg.Register("encoding/hex", p)
}

var _ = adt.TopKind // in case the adt package isn't used

var p = &pkg.Package{
	Native: []*pkg.Builtin{{
		Name: "EncodedLen",
		Params: []pkg.Param{
			{Kind: adt.IntKind},
		},
		Result: adt.IntKind,
		Func: func(c *pkg.CallCtxt) {
			n := c.Int(0)
			if c.Do() {
				c.Ret = EncodedLen(n)
			}
		},
	}, {
		Name: "DecodedLen",
		Params: []pkg.Param{
			{Kind: adt.IntKind},
		},
		Result: adt.IntKind,
		Func: func(c *pkg.CallCtxt) {
			x := c.Int(0)
			if c.Do() {
				c.Ret = DecodedLen(x)
			}
		},
	}, {
		Name: "Decode",
		Params: []pkg.Param{
			{Kind: adt.StringKind},
		},
		Result: adt.BytesKind | adt.StringKind,
		Func: func(c *pkg.CallCtxt) {
			s := c.String(0)
			if c.Do() {
				c.Ret, c.Err = Decode(s)
			}
		},
	}, {
		Name: "Dump",
		Params: []pkg.Param{
			{Kind: adt.BytesKind | adt.StringKind},
		},
		Result: adt.StringKind,
		Func: func(c *pkg.CallCtxt) {
			data := c.Bytes(0)
			if c.Do() {
				c.Ret = Dump(data)
			}
		},
	}, {
		Name: "Encode",
		Params: []pkg.Param{
			{Kind: adt.BytesKind | adt.StringKind},
		},
		Result: adt.StringKind,
		Func: func(c *pkg.CallCtxt) {
			src := c.Bytes(0)
			if c.Do() {
				c.Ret = Encode(src)
			}
		},
	}},
}
