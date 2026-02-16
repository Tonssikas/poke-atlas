import Bug from './bug.svg'
import Dark from './dark.svg'
import Dragon from './dragon.svg'
import Electric from './electric.svg'
import Fairy from './fairy.svg'
import Fighting from './fighting.svg'
import Fire from './fire.svg'
import Flying from './flying.svg'
import Ghost from './ghost.svg'
import Grass from './grass.svg'
import Ground from './ground.svg'
import Ice from './ice.svg'
import Normal from './normal.svg'
import Poison from './poison.svg'
import Psychic from './psychic.svg'
import Rock from './rock.svg'
import Steel from './steel.svg'
import Water from './water.svg'

export const typeIcons = {
    bug: Bug,
    dark: Dark,
    dragon: Dragon,
    electric: Electric,
    fairy: Fairy,
    fighting: Fighting,
    fire: Fire,
    flying: Flying,
    ghost: Ghost,
    grass: Grass,
    ground: Ground,
    ice: Ice,
    normal: Normal,
    poison: Poison,
    psychic: Psychic,
    rock: Rock,
    steel: Steel,
    water: Water
} as const;

export type PokemonType = keyof typeof typeIcons;