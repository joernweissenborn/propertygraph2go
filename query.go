package propertygraph2go

type GraphQuery struct{}

func Query() (q GraphQuery) { return }

func (q GraphQuery) Vertex(key Key) (vq VertexQuery) { return VertexQuery{v: q.vertex(key)} }
func (q GraphQuery) vertex(k Key) func(Graph) (Vertex, error) {
	return func(g Graph) (v Vertex, err error) { return g.GetVertex(k) }
}

func (q GraphQuery) Vertices() (vq VerticesQuery) { return VerticesQuery{v: q.vertices()} }
func (q GraphQuery) vertices() func(Graph) ([]Vertex, error) {
	return func(g Graph) (vs []Vertex, err error) {
		walkFn := func(v Vertex) {
			vs = append(vs, v)
		}
		g.WalkVertices(walkFn)
		return
	}
}

func (q GraphQuery) Edge(key Key) (vq EdgeQuery) { return EdgeQuery{q.edge(key)} }
func (q GraphQuery) edge(k Key) func(Graph) (Edge, error) {
	return func(g Graph) (e Edge, err error) { return g.GetEdge(k) }
}

type VertexQuery struct {
	v func(Graph) (Vertex, error)
}

func (q VertexQuery) Execute(g Graph) (v Vertex, err error) { return q.v(g) }

func (q VertexQuery) Incoming() EdgesQuery { return EdgesQuery{e: q.incoming} }
func (q VertexQuery) incoming(g Graph) (e []Edge, err error) {
	v, err := q.v(g)
	if err != nil {
		return
	}
	e = v.Incoming()
	return
}

func (q VertexQuery) Outgoing() EdgesQuery { return EdgesQuery{e: q.outgoing} }
func (q VertexQuery) outgoing(g Graph) (e []Edge, err error) {
	v, err := q.v(g)
	if err != nil {
		return
	}
	e = v.Outgoing()
	return
}

type VerticesQuery struct {
	v func(Graph) ([]Vertex, error)
}

type VertexFilter func(Vertex) bool

func (q VerticesQuery) Execute(g Graph) (v []Vertex, err error) { return q.v(g) }

func (q VerticesQuery) Where(filter ...VertexFilter) (v VerticesQuery) {
	return VerticesQuery{q.where(filter...)}
}
func (q VerticesQuery) where(filter ...VertexFilter) func(Graph) ([]Vertex, error) {
	return func(g Graph) (v []Vertex, err error) {
		vs, err := q.v(g)
		if err != nil {
			return
		}
		for _, cv := range vs {
			for _, f := range filter {
				if !f(cv) {
					goto CONTINUE
				}
			}
			v = append(v, cv)
		CONTINUE:
		}
		return
	}
}

func (q VerticesQuery) WhereAny(filter ...VertexFilter) (v VerticesQuery) {
	return VerticesQuery{q.whereany(filter...)}
}
func (q VerticesQuery) whereany(filter ...VertexFilter) func(Graph) ([]Vertex, error) {
	return func(g Graph) (v []Vertex, err error) {
		vs, err := q.v(g)
		if err != nil {
			return
		}
		for _, cv := range vs {
			for _, f := range filter {
				if f(cv) {
					v = append(v, cv)
					continue
				}
			}
		}
		return
	}
}

func (q VerticesQuery) HasProperty(key ...string) VerticesQuery {
	hasProp := func(k string) VertexFilter {
		return func(v Vertex) (has bool) {
			has = HasProperty(v, k)
			return
		}
	}
	f := []VertexFilter{}
	for _, k := range key {
		f = append(f, hasProp(k))
	}
	return q.Where(f...)
}

func (q VerticesQuery) HasAnyProperty(key ...string) VerticesQuery {
	hasProp := func(k string) VertexFilter {
		return func(v Vertex) (has bool) {
			has = HasProperty(v, k)
			return
		}
	}
	f := []VertexFilter{}
	for _, k := range key {
		f = append(f, hasProp(k))
	}
	return q.WhereAny(f...)
}

func (q VerticesQuery) HasPropertyValue(key string, value interface{}) VerticesQuery {
	hasProp := func(k string, val interface{}) VertexFilter {
		return func(v Vertex) (has bool) {
			has = HasPropertyValue(v, k, val)
			return
		}
	}
	return q.Where(hasProp(key, value))
}

func (q VerticesQuery) HasAnyPropertyValue(key string, value ...interface{}) VerticesQuery {
	hasProp := func(k string, val interface{}) VertexFilter {
		return func(v Vertex) (has bool) {
			has = HasPropertyValue(v, k, val)
			return
		}
	}
	f := []VertexFilter{}
	for _, v := range value {
		f = append(f, hasProp(key, v))
	}
	return q.WhereAny(f...)
}

func (q VerticesQuery) Incoming() EdgesQuery { return EdgesQuery{q.incoming} }
func (q VerticesQuery) incoming(g Graph) (e []Edge, err error) {
	vs, err := q.v(g)
	if err != nil {
		return
	}
	for _, v := range vs {
		e = append(e, v.Incoming()...)
	}
	return
}

func (q VerticesQuery) Outgoing() EdgesQuery { return EdgesQuery{q.outgoing} }
func (q VerticesQuery) outgoing(g Graph) (e []Edge, err error) {
	vs, err := q.v(g)
	if err != nil {
		return
	}
	for _, v := range vs {
		e = append(e, v.Outgoing()...)
	}
	return
}

type EdgeQuery struct {
	e func(Graph) (Edge, error)
}

func (q EdgeQuery) Execute(g Graph) (e Edge, err error) { return q.e(g) }

func (q EdgeQuery) Head() (v VertexQuery) { return VertexQuery{q.head} }
func (q EdgeQuery) head(g Graph) (v Vertex, err error) {
	e, err := q.e(g)
	if err != nil {
		return
	}
	v = e.Head()
	return
}

func (q EdgeQuery) Tail() (v VertexQuery) { return VertexQuery{q.tail} }
func (q EdgeQuery) tail(g Graph) (v Vertex, err error) {
	e, err := q.e(g)
	if err != nil {
		return
	}
	v = e.Tail()
	return
}

type EdgesQuery struct {
	e func(Graph) ([]Edge, error)
}

func (q EdgesQuery) Execute(g Graph) (e []Edge, err error) { return q.e(g) }

func (q EdgesQuery) Heads() (v VerticesQuery) { return VerticesQuery{q.heads} }
func (q EdgesQuery) heads(g Graph) (v []Vertex, err error) {
	es, err := q.e(g)
	if err != nil {
		return
	}
	seen := map[interface{}]interface{}{}
	for _, e := range es {
		if _, ok := seen[e.Head()]; ok {
			continue
		}
		seen[e.Head()] = nil
		v = append(v, e.Head())
	}
	return
}

func (q EdgesQuery) Tails() (v VerticesQuery) { return VerticesQuery{q.tails} }
func (q EdgesQuery) tails(g Graph) (v []Vertex, err error) {
	es, err := q.e(g)
	if err != nil {
		return
	}
	seen := map[interface{}]interface{}{}
	for _, e := range es {
		if _, ok := seen[e.Tail()]; ok {
			continue
		}
		seen[e.Tail()] = nil
		v = append(v, e.Tail())
	}
	return
}

type EdgeFilter func(Edge) bool

func (q EdgesQuery) Where(filter ...EdgeFilter) (v EdgesQuery) {
	return EdgesQuery{q.where(filter...)}
}
func (q EdgesQuery) where(filter ...EdgeFilter) func(Graph) ([]Edge, error) {
	return func(g Graph) (v []Edge, err error) {
		es, err := q.e(g)
		if err != nil {
			return
		}
		for _, ce := range es {
			for _, f := range filter {
				if !f(ce) {
					goto CONTINUE
				}
			}
			v = append(v, ce)
		CONTINUE:
		}
		return
	}
}

func (q EdgesQuery) WhereAny(filter ...EdgeFilter) (v EdgesQuery) {
	return EdgesQuery{q.whereAny(filter...)}
}
func (q EdgesQuery) whereAny(filter ...EdgeFilter) func(Graph) ([]Edge, error) {
	return func(g Graph) (v []Edge, err error) {
		es, err := q.e(g)
		if err != nil {
			return
		}
		for _, ce := range es {
			for _, f := range filter {
				if f(ce) {
			v = append(v, ce)
			continue
				}
			}
		}
		return
	}
}

func (q EdgesQuery) HasLabel(label ...string) EdgesQuery {
	hasLabel := func(lbl string) EdgeFilter {
		return func(e Edge) (has bool) { return lbl == e.Label() }
	}
	lblFs := []EdgeFilter{}
	for _, l := range label {
		lblFs = append(lblFs, hasLabel(l))
	}
	return q.Where(lblFs...)
}

func (q EdgesQuery) HasProperty(key ...string) EdgesQuery {
	hasProp := func(k string) EdgeFilter {
		return func(e Edge) (has bool) { return HasProperty(e, k) }
	}
	f := []EdgeFilter{}
	for _, k := range key {
		f = append(f, hasProp(k))
	}
	return q.Where(f...)
}

func (q EdgesQuery) HasAnyProperty(key ...string) EdgesQuery {
	hasProp := func(k string) EdgeFilter {
		return func(e Edge) (has bool) { return HasProperty(e, k) }
	}
	f := []EdgeFilter{}
	for _, k := range key {
		f = append(f, hasProp(k))
	}
	return q.WhereAny(f...)
}

func (q EdgesQuery) HasPropertyValue(key string, value interface{}) EdgesQuery {
	hasProp := func(k string, v interface{}) EdgeFilter {
		return func(e Edge) (has bool) { return HasPropertyValue(e, k, v) }
	}
	return q.Where(hasProp(key, value))
}

func (q EdgesQuery) HasAnyPropertyValue(key string, value ...interface{}) EdgesQuery {
	hasProp := func(k string, v interface{}) EdgeFilter {
		return func(e Edge) (has bool) { return HasPropertyValue(e, k, v) }
	}
	f := []EdgeFilter{}
	for _, v := range value {
		f = append(f, hasProp(key, v))
	}
	return q.WhereAny(f...)
}

func (q EdgesQuery) HeadHasProperty(key ...string) EdgesQuery {
	hasProp := func(k string) EdgeFilter {
		return func(e Edge) (has bool) { return HasProperty(e.Head(), k) }
	}
	f := []EdgeFilter{}
	for _, k := range key {
		f = append(f, hasProp(k))
	}
	return q.Where(f...)
}

func (q EdgesQuery) HeadHasAnyProperty(key ...string) EdgesQuery {
	hasProp := func(k string) EdgeFilter {
		return func(e Edge) (has bool) {
			has = HasProperty(e.Head(), k)
			return
		}
	}
	f := []EdgeFilter{}
	for _, k := range key {
		f = append(f, hasProp(k))
	}
	return q.WhereAny(f...)
}

func (q EdgesQuery) HeadHasPropertyValue(key string, value interface{}) EdgesQuery {
	hasProp := func(k string, v interface{}) EdgeFilter {
		return func(e Edge) (has bool) {
			has = HasPropertyValue(e.Head(), k, v)
			return
		}
	}
	return q.Where(hasProp(key, value))
}

func (q EdgesQuery) HeadHasAnyPropertyValue(key string, value ...interface{}) EdgesQuery {
	hasProp := func(k string, v interface{}) EdgeFilter {
		return func(e Edge) (has bool) {return HasPropertyValue(e.Head(), k, v)		}
	}
	f := []EdgeFilter{}
	for _, v := range value {
		f = append(f, hasProp(key,v))
	}
	return q.WhereAny(f...)
}

func (q EdgesQuery) TailHasProperty(key ...string) EdgesQuery {
	hasProp := func(k string) EdgeFilter {
		return func(e Edge) (has bool) {	return HasProperty(e.Tail(), k)		}
	}
	f := []EdgeFilter{}
	for _, k := range key {
		f = append(f, hasProp(k))
	}
	return q.Where(f...)
}

func (q EdgesQuery) TailHasAnyProperty(key ...string) EdgesQuery {
	hasProp := func(k string) EdgeFilter {
		return func(e Edge) (has bool) {	return HasProperty(e.Tail(), k)		}
	}
	f := []EdgeFilter{}
	for _, k := range key {
		f = append(f, hasProp(k))
	}
	return q.WhereAny(f...)
}

func (q EdgesQuery) TailHasPropertyValue(key string, value interface{}) EdgesQuery {
	hasProp := func(k string, v interface{}) EdgeFilter {
		return func(e Edge) (has bool) {return HasPropertyValue(e.Tail(), k, v)		}
	}
	return q.Where(hasProp(key, value))
}
func (q EdgesQuery) TailAnyPropertyValue(key string, value ...interface{}) EdgesQuery {
	hasProp := func(k string, v interface{}) EdgeFilter {
		return func(e Edge) (has bool) {return HasPropertyValue(e.Tail(), k, v)		}
	}
	f := []EdgeFilter{}
	for _, v := range value {
		f = append(f, hasProp(key,v))
	}
	return q.WhereAny(f...)
}
