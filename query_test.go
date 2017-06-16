package propertygraph2go

import "testing"

	var kg = NewKeyGen()

func testGraph() (g Graph) {
	g = NewSimpleGraph()

	rankEnsign, _ := g.CreateVertex(kg.NextVertex())
	rankEnsign.SetProperty("rank", "ensign")
	rankLieutenant, _ := g.CreateVertex(kg.NextVertex())
	rankLieutenant.SetProperty("rank", "lieutenant")
	g.CreateEdge(kg.NextEdge(), "outranks", rankLieutenant.Key(), rankEnsign.Key())
	rankLtCommander, _ := g.CreateVertex(kg.NextVertex())
	rankLtCommander.SetProperty("rank", "lieutenant commander")
	g.CreateEdge(kg.NextEdge(), "outranks", rankLtCommander.Key(), rankEnsign.Key())
	g.CreateEdge(kg.NextEdge(), "outranks", rankLtCommander.Key(), rankLieutenant.Key())
	rankCommander, _ := g.CreateVertex(kg.NextVertex())
	rankCommander.SetProperty("rank", "commander")
	g.CreateEdge(kg.NextEdge(), "outranks", rankCommander.Key(), rankEnsign.Key())
	g.CreateEdge(kg.NextEdge(), "outranks", rankCommander.Key(), rankLieutenant.Key())
	g.CreateEdge(kg.NextEdge(), "outranks", rankCommander.Key(), rankLtCommander.Key())
	rankCaptain, _ := g.CreateVertex(kg.NextVertex())
	rankCaptain.SetProperty("rank", "captain")
	g.CreateEdge(kg.NextEdge(), "outranks", rankCaptain.Key(), rankEnsign.Key())
	g.CreateEdge(kg.NextEdge(), "outranks", rankCaptain.Key(), rankLieutenant.Key())
	g.CreateEdge(kg.NextEdge(), "outranks", rankCaptain.Key(), rankLtCommander.Key())
	g.CreateEdge(kg.NextEdge(), "outranks", rankCaptain.Key(), rankCommander.Key())

	posCouns, _ := g.CreateVertex(kg.NextVertex())
	posCouns.SetProperty("position", "counselor")
	posCom, _ := g.CreateVertex(kg.NextVertex())
	posCom.SetProperty("position", "communication")
	posAstro, _ := g.CreateVertex(kg.NextVertex())
	posAstro.SetProperty("position", "astronomy")
	posHelm, _ := g.CreateVertex(kg.NextVertex())
	posHelm.SetProperty("position", "helm")
	posSecurity, _ := g.CreateVertex(kg.NextVertex())
	posSecurity.SetProperty("position", "security")
	posScience, _ := g.CreateVertex(kg.NextVertex())
	posScience.SetProperty("position", "science")
	posEngineering, _ := g.CreateVertex(kg.NextVertex())
	posEngineering.SetProperty("position", "engineering")
	posMedical, _ := g.CreateVertex(kg.NextVertex())
	posMedical.SetProperty("position", "medical")
	posFirstOff, _ := g.CreateVertex(kg.NextVertex())
	posFirstOff.SetProperty("position", "first officer")
	posCmdOff, _ := g.CreateVertex(kg.NextVertex())
	posCmdOff.SetProperty("position", "commanding officer")

	crewmember := func(g Graph, ship Vertex, name, firstname string) (m Vertex) {
		m, _ = g.CreateVertex(kg.NextVertex())
		g.CreateEdge(kg.NextEdge(), "crewmember", m.Key(), ship.Key())

		m.SetProperty("name", name)
		m.SetProperty("firstname", firstname)
		return
	}

	enterpriseA, _ := g.CreateVertex(kg.NextVertex())
	enterpriseA.SetProperty("ship", nil)
	enterpriseA.SetProperty("name", "enterprise")
	enterpriseA.SetProperty("register", "NCC-1701-A")

	c := crewmember(g, enterpriseA, "kirk", "jim").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankCaptain.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posCmdOff.Key(), c)

	c = crewmember(g, enterpriseA, "spock", "").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankCommander.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posFirstOff.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posScience.Key(), c)

	c = crewmember(g, enterpriseA, "mccoy", "leonard").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLtCommander.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posMedical.Key(), c)

	c = crewmember(g, enterpriseA, "scott", "montgomerry").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLieutenant.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posEngineering.Key(), c)

	c = crewmember(g, enterpriseA, "sulu", "hikaro").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLieutenant.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posHelm.Key(), c)

	c = crewmember(g, enterpriseA, "chekov", "pavel").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankEnsign.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posSecurity.Key(), c)

	c = crewmember(g, enterpriseA, "uhura", "").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLieutenant.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posCom.Key(), c)
	// 7

	enterpriseD, _ := g.CreateVertex(kg.NextVertex())
	enterpriseD.SetProperty("ship", nil)
	enterpriseD.SetProperty("name", "enterprise")
	enterpriseD.SetProperty("register", "NCC-1701-D")

	c = crewmember(g, enterpriseD, "picard", "jean-luc").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankCaptain.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posCmdOff.Key(), c)

	c = crewmember(g, enterpriseD, "riker", "william").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankCommander.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posFirstOff.Key(), c)

	c = crewmember(g, enterpriseD, "crusher", "beverly").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankCommander.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posMedical.Key(), c)

	c = crewmember(g, enterpriseD, "data", "").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLtCommander.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posScience.Key(), c)

	c = crewmember(g, enterpriseD, "worf", "").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLieutenant.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posSecurity.Key(), c)

	c = crewmember(g, enterpriseD, "laforge", "geordi").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLieutenant.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posEngineering.Key(), c)

	c = crewmember(g, enterpriseD, "troy", "diana").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLtCommander.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posCouns.Key(), c)
	// 7

	voyager, _ := g.CreateVertex(kg.NextVertex())
	voyager.SetProperty("ship", nil)
	voyager.SetProperty("name", "voyager")
	voyager.SetProperty("register", "NCC-74656")

	c = crewmember(g, voyager, "janeway", "catherine").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankCaptain.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posCmdOff.Key(), c)

	c = crewmember(g, voyager, "chakotay", "").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankCommander.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posFirstOff.Key(), c)

	c = crewmember(g, voyager, "medical emergency hologramm", "").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankCommander.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posMedical.Key(), c)

	c = crewmember(g, voyager, "kim", "harry").Key()
	g.CreateEdge(kg.NextEdge(), "position", posScience.Key(), c)
	g.CreateEdge(kg.NextEdge(), "rank", rankEnsign.Key(), c)

	c = crewmember(g, voyager, "tuvok", "").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLtCommander.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posSecurity.Key(), c)

	c = crewmember(g, voyager, "torres", "b'lana").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLieutenant.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posEngineering.Key(), c)

	c = crewmember(g, voyager, "paris", "tom").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankLieutenant.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posHelm.Key(), c)

	c = crewmember(g, voyager, "seven of nine", "").Key()
	g.CreateEdge(kg.NextEdge(), "rank", rankEnsign.Key(), c)
	g.CreateEdge(kg.NextEdge(), "position", posAstro.Key(), c)
	// 8
	return
}

func TestHowManyCrewOnAllEnt(t *testing.T) {
	ships := Query().Vertices().HasProperty("ship")
	s, err := ships.Execute(testGraph())
	if err != nil {
		t.Error(err)
	}
	if len(s) != 3 {
		t.Error("Wrong Nr of ships", 3, len(s))
	}
	ents := ships.HasPropertyValue("name", "enterprise")
	e, err := ents.Execute(testGraph())
	if err != nil {
		t.Error(err)
	}
	if len(e) != 2 {
		t.Error("Wrong Nr of enterprises", 2, len(e))
	}
	q := ents.Outgoing().HasLabel("crewmember").Heads()
	crew, err := q.Execute(testGraph())
	if err != nil {
		t.Error(err)
	}
	if len(crew) != 14 {
		t.Error("Wrong Nr of CrewMember", 14, len(crew))
	}
}

func TestHowManyOutranksKim(t *testing.T) {
	q := Query().Vertices().HasPropertyValue("name", "kim").Outgoing().HasLabel("rank").Heads().Outgoing().HasLabel("outranks").Heads().Incoming().HasLabel("rank").Tails()
	outranks, err := q.Execute(testGraph())
	if err != nil {
		t.Error(err)
	}
	if len(outranks) != 19 {
		t.Error("Wrong Nr of CrewMember", 19, len(outranks))
	}
}

func TestHowManyLtsOnEnterpriseD(t *testing.T) {
	crew := Query().Vertices().HasProperty("ship").HasPropertyValue("register", "NCC-1701-D").Outgoing().HasLabel("crewmember").Heads()
	c, err := crew.Execute(testGraph())
	if err != nil {
		t.Error(err)
	}
	if len(c) != 7 {
		t.Error("Wrong Nr of CrewMember", 7, len(c))
	}
	q := crew.Outgoing().HeadHasPropertyValue("rank", "lieutenant").Tails()
	lts, err := q.Execute(testGraph())
	if err != nil {
		t.Error(err)
	}
	if len(lts) != 2 {
		t.Error("Wrong Nr of CrewMember", 2, len(lts))
	}
}

func TestHowManyLtCmdsAll(t *testing.T) {
	crew := Query().Vertices().HasProperty("ship").Outgoing().HasLabel("crewmember").Heads()
	c, err := crew.Execute(testGraph())
	if err != nil {
		t.Error(err)
	}
	if len(c) != 22 {
		t.Error("Wrong Nr of CrewMember", 22, len(c))
	}
	ltcmds := crew.Outgoing().HeadHasPropertyValue("rank", "lieutenant commander").Tails()
	lt, err := ltcmds.Execute(testGraph())
	if err != nil {
		t.Error(err)
	}
	if len(lt) != 4 {
		t.Error("Wrong Nr of LtCmds", 4, len(lt))
	}
}

func TestHowManyOneName(t *testing.T) {
	q := Query().Vertices().HasProperty("ship").Outgoing().HasLabel("crewmember").Heads().HasPropertyValue("firstname", "")
	crew, err := q.Execute(testGraph())
	if err != nil {
		t.Error(err)
	}
	if len(crew) != 8 {
		t.Error("Wrong Nr of CrewMember", 8, len(crew))
	}
}

func TestAddUniformColors(t *testing.T) {
	crewClassic := Query().Vertices().HasProperty("ship").HasPropertyValue("register", "NCC-1701-A").Outgoing().HasLabel("crewmember").Heads()
	crewModern := Query().Vertices().HasProperty("ship").HasAnyPropertyValue("register", "NCC-1701-D", "NCC-74656").Outgoing().HasLabel("crewmember").Heads()

	g := testGraph()
	redsClassic, err := crewClassic.Outgoing().HasLabel("position").HeadHasAnyPropertyValue("position", "security", "engineering", "communication").Tails().Execute(g)
	if err != nil {
		t.Error(err)
	}
	redsModern, err := crewModern.Outgoing().HasLabel("position").HeadHasAnyPropertyValue("position", "commanding officer", "first officer", "helm").Tails().Execute(g)
	if err != nil {
		t.Error(err)
	}

	red, err := g.CreateVertex(kg.NextVertex())
	if err != nil {
		t.Error(err)
	}
	err = red.SetProperty("uniform", "red")
	if err != nil {
		t.Error(err)
	}

	for _, c := range append(redsClassic, redsModern...) {
		_,err = g.CreateEdge(kg.NextEdge(),"uniform", red.Key(), c.Key())
		if err != nil {
			t.Error(err)
		}
	}

	reds, err := Query().Vertices().HasPropertyValue("uniform", "red").Incoming().HasLabel("uniform").Execute(g)
	if err != nil {
		t.Error(err)
	}
	if len(reds) != 8 {
		t.Error("Wrong Nr of CrewMember", 8, len(reds))
	}

}
