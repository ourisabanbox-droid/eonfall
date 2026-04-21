package world

type CivilizationAxis string

const (
	CivilizationAxisResilience CivilizationAxis = "resilience"
	CivilizationAxisExpansion  CivilizationAxis = "expansion"
	CivilizationAxisInfluence  CivilizationAxis = "influence"
	CivilizationAxisScience    CivilizationAxis = "science"
)

type CivilizationProfile struct {
	DominantAxis  CivilizationAxis
	SecondaryAxis CivilizationAxis
	ProfileLabel  string
}

func ComputeCivilizationProfile(t *CivilizationTrajectory) CivilizationProfile {
	scores := []struct {
		Axis  CivilizationAxis
		Score int
	}{
		{Axis: CivilizationAxisResilience, Score: t.ResilienceScore},
		{Axis: CivilizationAxisExpansion, Score: t.ExpansionScore},
		{Axis: CivilizationAxisInfluence, Score: t.InfluenceScore},
		{Axis: CivilizationAxisScience, Score: t.ScienceScore},
	}

	// tri simple décroissant par score
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].Score > scores[i].Score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	dominant := scores[0].Axis
	secondary := scores[1].Axis

	return CivilizationProfile{
		DominantAxis:  dominant,
		SecondaryAxis: secondary,
		ProfileLabel:  computeCivilizationProfileLabel(dominant, secondary),
	}
}

func computeCivilizationProfileLabel(dominant, secondary CivilizationAxis) string {
	switch dominant {
	case CivilizationAxisResilience:
		switch secondary {
		case CivilizationAxisInfluence:
			return "Administration de crise"
		case CivilizationAxisScience:
			return "Civilisation de redressement"
		case CivilizationAxisExpansion:
			return "Frontière résiliente"
		default:
			return "Culture de survie"
		}
	case CivilizationAxisScience:
		switch secondary {
		case CivilizationAxisResilience:
			return "Technocratie de survie"
		case CivilizationAxisInfluence:
			return "Rationalistes d'influence"
		case CivilizationAxisExpansion:
			return "Expansion scientifique"
		default:
			return "Technocratie émergente"
		}
	case CivilizationAxisInfluence:
		switch secondary {
		case CivilizationAxisResilience:
			return "Ordre conciliateur"
		case CivilizationAxisScience:
			return "Diplomatie technicienne"
		case CivilizationAxisExpansion:
			return "Hégémonie politique"
		default:
			return "Puissance d'influence"
		}
	case CivilizationAxisExpansion:
		switch secondary {
		case CivilizationAxisResilience:
			return "Empire tenace"
		case CivilizationAxisScience:
			return "Expansion planifiée"
		case CivilizationAxisInfluence:
			return "Empire des pactes"
		default:
			return "Ambition impériale"
		}
	default:
		return "Civilisation émergente"
	}
}
