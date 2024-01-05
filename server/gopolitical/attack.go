package gopolitical

type AttackAction struct {
	Action
	AttackerCountry   *Country
	DefenderTerritory *Territory
	ArmamentUsed      float64
}

func (a *AttackAction) Execute(env *Environment) {
	// On vérifie que les pays sont bien voisins au moments de l'attaques
	if !env.World.IsTerritoryAttackableBy(a.DefenderTerritory, a.AttackerCountry) {
		Debug(
			a.AttackerCountry.Name,
			"[Environment] Attaque avortée sur %s (%s) pour cause de cible inatteignable\n",
			a.AttackerCountry.Name,
			a.DefenderTerritory.Country.Name,
		)
		return
	}
	// On récupère les stocks d'armes des deux pays
	defensiveArmament := a.DefenderTerritory.Country.GetTotalStockOf(ARMAMENT)
	offensiveArmament := a.AttackerCountry.GetTotalStockOf(ARMAMENT)
	// On vérifie que le pays à bien de quoi attaquer
	if offensiveArmament < a.ArmamentUsed && a.ArmamentUsed > 0 && offensiveArmament > 0 {
		Debug(
			a.AttackerCountry.Name,
			"[Environment] Attaque avortée sur %s (%s) pour cause de manque d'armement\n",
			a.AttackerCountry.Name,
			a.DefenderTerritory.Country.Name,
		)
		return
	}
	// Il n'utilisera que ce qu'il souhaite
	offensiveArmament = a.ArmamentUsed
	// Si l'attaqué a assez de ressource pour se défendre
	if defensiveArmament < 0 {
		defensiveArmament = 0
	}
	if a.ArmamentUsed < defensiveArmament {
		defensiveArmament = a.ArmamentUsed // Il n'en utilise qu'une partie
	}

	// Il font une bataille donc il consomme de l'armement
	a.AttackerCountry.Consume(ARMAMENT, offensiveArmament)
	a.DefenderTerritory.Country.Consume(ARMAMENT, defensiveArmament)

	// Le taux de réussite correspond à offensif / défensif avec un bonus de
	chanceOfCapture := 1 - (offensiveArmament/defensiveArmament)*(1-INCOMPRESSIBLE_CAPTURE_RATIO)
	Debug(
		a.AttackerCountry.Name,
		"Attaque %v (%v) avec %.0f%% de réussite",
		a.DefenderTerritory.Name,
		a.DefenderTerritory.Country.Name,
		chanceOfCapture*100,
	)

	// On récupère l'état de la relation actuelle
	relation := env.RelationManager.GetRelation(a.AttackerCountry.ID, a.DefenderTerritory.Country.ID)
	attackedCountry := a.DefenderTerritory.Country
	if chanceOfCapture > random.Float64() { // L'attaque a réussi
		relation = relation * RELATION_RATIO_ATTACK
		a.DefenderTerritory.TransfertProperty(a.AttackerCountry)
		Debug("Environment", "Capturé !")
	} else { // L'attaque a échoué
		relation = relation * RELATION_RATIO_DEFEND
		Debug("Environment", "Échec !")
	}
	env.RelationManager.UpdateRelation(a.AttackerCountry.ID, attackedCountry.ID, relation)
}
