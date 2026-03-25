import type { PokemonType } from "$lib/icons/pokemon-types";

type TypeRelations = {
  strong: PokemonType[];
  weak: PokemonType[];
  immune: PokemonType[];
};

const attackRelations: Record<PokemonType, TypeRelations> = {
    fire: {
        strong: ["grass", "ice", "bug", "steel"],
        weak: ["fire", "water", "rock"],
        immune: []
    },
    water: {
        strong: ["fire", "ground", "rock"],
        weak: ["water", "grass", "dragon"],
        immune: []
    },
    bug: {
        strong: ["grass", "psychic", "dark"],
        weak: ["fire", "fighting", "poison", "flying", "ghost", "steel", "fairy"],
        immune: []
    },
    dark: {
        strong: ["psychic", "ghost"],
        weak: ["fighting", "dark", "fairy"],
        immune: []
    },
    dragon: {
        strong: ["dragon"],
        weak: ["steel"],
        immune: ["fairy"]
    },
    electric: {
        strong: ["water", "flying"],
        weak: ["electric", "grass", "dragon"],
        immune: ["ground"]
    },
    fairy: {
        strong: ["fighting", "dragon", "dark"],
        weak: ["fire", "poison", "steel"],
        immune: []
    },
    fighting: {
        strong: ["normal", "ice", "rock", "dark", "steel"],
        weak: ["poison", "flying", "psychic", "bug", "fairy"],
        immune: ["ghost"]
    },
    flying: {
        strong: ["grass", "fighting", "bug"],
        weak: ["electric", "rock", "steel"],
        immune: []
    },
    ghost: {
        strong: ["psychic", "ghost"],
        weak: ["dark"],
        immune: ["normal"]
    },
    grass: {
        strong: ["water", "ground", "rock"],
        weak: ["fire", "grass", "poison", "flying", "bug", "dragon", "steel"],
        immune: []
    },
    ground: {
        strong: ["fire", "electric", "poison", "rock", "steel"],
        weak: ["grass", "bug"],
        immune: ["flying"]
    },
    ice: {
        strong: ["grass", "ground", "flying", "dragon"],
        weak: ["fire", "water", "ice", "steel"],
        immune: []
    },
    normal: {
        strong: [],
        weak: ["rock", "steel"],
        immune: ["ghost"]
    },
    poison: {
        strong: ["grass", "fairy"],
        weak: ["poison", "ground", "rock", "ghost"],
        immune: ["steel"]
    },
    psychic: {
        strong: ["fighting", "poison"],
        weak: ["psychic", "steel"],
        immune: ["dark"]
    },
    rock: {
        strong: ["fire", "ice", "flying", "bug"],
        weak: ["fighting", "ground", "steel"],
        immune: []
    },
    steel: {
        strong: ["ice", "rock", "fairy"],
        weak: ["fire", "water", "electric", "steel"],
        immune: []
    },
}

function getDefensiveMultiplier(
  defendingTypes: PokemonType[],
  attackingType: PokemonType
): number {
  return defendingTypes.reduce((total, defType) => {
    const rel = attackRelations[attackingType];

    if (rel.immune.includes(defType)) return 0;
    if (rel.strong.includes(defType)) return total * 2;
    if (rel.weak.includes(defType)) return total * 0.5;

    return total;
  }, 1);
}

function getTypeEffectiveness(defendingTypes: PokemonType[]) {
    const result: Record<PokemonType, number> = {} as any;

    for (const attackType in attackRelations) {
        const multiplier = getDefensiveMultiplier(defendingTypes, attackType as PokemonType);

        result[attackType as PokemonType] = multiplier;
    }
    return result;
}

export function groupEffectiveness(defendingTypes: PokemonType[]) {
    const effectiveness = getTypeEffectiveness(defendingTypes);

    const groups = {
        x4: [] as PokemonType[],
        x2: [] as PokemonType[],
        x05: [] as PokemonType[],
        x025: [] as PokemonType[],
        immune: [] as PokemonType[],
    };

    for (const type in effectiveness) {
        const value = effectiveness[type as PokemonType];

        if (value === 4) groups.x4.push(type as PokemonType);
        else if (value === 2) groups.x2.push(type as PokemonType);
        else if (value === 0.5) groups.x05.push(type as PokemonType);
        else if (value === 0.25) groups.x025.push(type as PokemonType);
        else if (value === 0) groups.immune.push(type as PokemonType);
    }

    return groups;
}