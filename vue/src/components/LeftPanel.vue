<template>
    <div class="p-4 text-left w-full">
        <h1 class="text-xl font-bold mb-4">Категории</h1>
        <TreeCategory :json="data" :expanded-paths="expandedPaths" @toggle="toggleExpanded" @selectItem="selectItem"
            />
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
 
import TreeCategory from './TreeCategory.vue';

const expandedPaths = ref<Set<string>>(new Set());

const data = {
    "Condenser Packs": {
        "Advanced Condenser Packs": [
            "ElectroPunch Ultra L",
            "StrikeSnipe Ultra L",

            "ElectroPunch Ultra M",
            "StrikeSnipe Ultra M",

            "ElectroPunch Ultra S",
            "StrikeSnipe Ultra S"
        ],
        "Standard Condenser Packs": [
            "BlastShot Condenser Pack L",
            "GalvaSurge Condenser Pack L",
            "MesmerFlux Condenser Pack L",
            "SlamBolt Condenser Pack L",

            "BlastShot Condenser Pack M",
            "GalvaSurge Condenser Pack M",
            "MesmerFlux Condenser Pack M",
            "SlamBolt Condenser Pack M",

            "BlastShot Condenser Pack S",
            "GalvaSurge Condenser Pack S",
            "MesmerFlux Condenser Pack S",
            "SlamBolt Condenser Pack S"
        ]
    },
    "Bombs": [
        "Concussion Bomb",
        "Electron Bomb",
        "Focused Void Bomb",
        "Lockbreaker Bomb",
        "Scorch Bomb",
        "Shrapnel Bomb",
        "Void Bomb"
    ],
    "Breacher Pods": ["SCARAB Breacher Pod M", "SCARAB Breacher Pod S"],

    "Missiles": {
        "Auto-Targeting": [
            "Inferno Auto-Targeting Heavy Missile I",
            "Scourge Auto-Targeting Cruise Missile I",
            "Mjolnir Auto-Targeting Cruise Missile I",
            "Nova Auto-Targeting Cruise Missile I",
            "Inferno Auto-Targeting Cruise Missile I",
            "Mjolnir Auto-Targeting Light Missile I",
            "Scourge Auto-Targeting Light Missile I",
            "Nova Auto-Targeting Light Missile I",
            "Inferno Auto-Targeting Light Missile I",
            "Scourge Auto-Targeting Heavy Missile I",
            "Mjolnir Auto-Targeting Heavy Missile I",
            "Nova Auto-Targeting Heavy Missile I",

            "Legion Scourge Auto-Targeting Heavy Missile",
            "Legion Mjolnir Auto-Targeting Cruise Missile",
            "Legion Nova Auto-Targeting Light Missile",
            "Legion Inferno Auto-Targeting Heavy Missile",
            "Legion Scourge Auto-Targeting Cruise Missile",
            "Legion Nova Auto-Targeting Heavy Missile",
            "Legion Inferno Auto-Targeting Cruise Missile",
            "Legion Mjolnir Auto-Targeting Light Missile",
            "Legion Nova Auto-Targeting Cruise Missile",
            "Legion Scourge Auto-Targeting Light Missile",
            "Legion Mjolnir Auto-Targeting Heavy Missile",
            "Legion Inferno Auto-Targeting Light Missile"
        ],
        "Defender": [
            "Defender Missile I"
        ],
        "Light Missiles": [
            "Inferno Fury Light Missile",
            "Mjolnir Fury Light Missile",
            "Nova Fury Light Missile",
            "Scourge Fury Light Missile",

            "Inferno Precision Light Missile",
            "Mjolnir Precision Light Missile",
            "Nova Precision Light Missile",
            "Scourge Precision Light Missile",

            "Caldari Navy Scourge Light Missile",
            "Dread Guristas Scourge Light Missile",
            "Guristas Inferno Light Missile",
            "Dread Guristas Inferno Light Missile",
            "Guristas Scourge Light Missile",
            "Caldari Navy Inferno Light Missile",
            "Dread Guristas Nova Light Missile",
            "Guristas Nova Light Missile",
            "Caldari Navy Nova Light Missile",
            "Guristas Mjolnir Light Missile",
            "Dread Guristas Mjolnir Light Missile",
            "Caldari Navy Mjolnir Light Missile",

            "Scourge Light Missile",
            "Inferno Light Missile",
            "Mjolnir Light Missile",
            "Nova Light Missile"
        ],
        "Rockets": [
            "Mjolnir Rocket",
            "Scourge Rocket",
            "Nova Rocket",
            "Inferno Rocket",

            "Mjolnir Javelin Rocket",
            "Scourge Javelin Rocket",
            "Nova Javelin Rocket",
            "Inferno Javelin Rocket",

            "Nova Rage Rocket",
            "Inferno Rage Rocket",
            "Scourge Rage Rocket",
            "Mjolnir Rage Rocket",

            "Dread Guristas Nova Rocket",
            "Guristas Scourge Rocket",
            "Caldari Navy Scourge Rocket",
            "Dread Guristas Scourge Rocket",
            "Guristas Inferno Rocket",
            "Caldari Navy Inferno Rocket",
            "Dread Guristas Inferno Rocket",
            "Guristas Mjolnir Rocket",
            "Caldari Navy Mjolnir Rocket",
            "Dread Guristas Mjolnir Rocket",
            "Guristas Nova Rocket",
            "Caldari Navy Nova Rocket"
        ],
        "Torpedoes": [
            "Mjolnir Torpedo",
            "Scourge Torpedo",
            "Nova Torpedo",
            "Inferno Torpedo",

            "Nova Javelin Torpedo",
            "Mjolnir Javelin Torpedo",
            "Inferno Javelin Torpedo",
            "Scourge Javelin Torpedo",

            "Scourge Rage Torpedo",
            "Inferno Rage Torpedo",
            "Mjolnir Rage Torpedo",
            "Nova Rage Torpedo",

            "Dread Guristas Nova Torpedo",
            "Guristas Mjolnir Torpedo",
            "Caldari Navy Mjolnir Torpedo",
            "Dread Guristas Mjolnir Torpedo",
            "Guristas Scourge Torpedo",
            "Caldari Navy Scourge Torpedo",
            "Dread Guristas Scourge Torpedo",
            "Guristas Inferno Torpedo",
            "Caldari Navy Inferno Torpedo",
            "Dread Guristas Inferno Torpedo",
            "Guristas Nova Torpedo",
            "Caldari Navy Nova Torpedo"
        ],

        "XL Torpedoes": [
            "Mjolnir XL Torpedo",
            "Scourge XL Torpedo",
            "Inferno XL Torpedo",
            "Nova XL Torpedo",

            "Guristas Nova XL Torpedo",
            "Guristas Inferno XL Torpedo",
            "Guristas Scourge XL Torpedo",
            "Guristas Mjolnir XL Torpedo",

            "Nova Rage XL Torpedo",
            "Mjolnir Rage XL Torpedo",
            "Inferno Rage XL Torpedo",
            "Scourge Rage XL Torpedo",

            "Nova Javelin XL Torpedo",
            "Mjolnir Javelin XL Torpedo",
            "Inferno Javelin XL Torpedo",
            "Scourge Javelin XL Torpedo"
        ],

        "Cruise Missiles": [
            "Nova Precision Cruise Missile",
            "Mjolnir Precision Cruise Missile",
            "Inferno Precision Cruise Missile",
            "Scourge Precision Cruise Missile",

            "Mjolnir Cruise Missile",
            "Scourge Cruise Missile",
            "Inferno Cruise Missile",
            "Nova Cruise Missile",

            "Scourge Fury Cruise Missile",
            "Nova Fury Cruise Missile",
            "Inferno Fury Cruise Missile",
            "Mjolnir Fury Cruise Missile",

            "Dread Guristas Inferno Cruise Missile",
            "Caldari Navy Scourge Cruise Missile",
            "Caldari Navy Nova Cruise Missile",
            "Dread Guristas Scourge Cruise Missile",
            "Dread Guristas Nova Cruise Missile",
            "Guristas Mjolnir Cruise Missile",
            "Caldari Navy Mjolnir Cruise Missile",
            "Guristas Nova Cruise Missile",
            "Guristas Scourge Cruise Missile",
            "Guristas Inferno Cruise Missile",
            "Dread Guristas Mjolnir Cruise Missile",
            "Caldari Navy Inferno Cruise Missile"
        ],

        "Heavy Missiles": [
            "Scourge Precision Heavy Missile",
            "Inferno Precision Heavy Missile",
            "Mjolnir Precision Heavy Missile",
            "Nova Precision Heavy Missile",

            "Inferno Heavy Missile",
            "Scourge Heavy Missile",
            "Nova Heavy Missile",
            "Mjolnir Heavy Missile",

            "Mjolnir Fury Heavy Missile",
            "Nova Fury Heavy Missile",
            "Scourge Fury Heavy Missile",
            "Inferno Fury Heavy Missile",

            "Guristas Mjolnir Heavy Missile",
            "Caldari Navy Mjolnir Heavy Missile",
            "Dread Guristas Mjolnir Heavy Missile",
            "Guristas Scourge Heavy Missile",
            "Caldari Navy Scourge Heavy Missile",
            "Dread Guristas Scourge Heavy Missile",
            "Guristas Inferno Heavy Missile",
            "Caldari Navy Inferno Heavy Missile",
            "Dread Guristas Inferno Heavy Missile",
            "Guristas Nova Heavy Missile",
            "Caldari Navy Nova Heavy Missile",
            "Dread Guristas Nova Heavy Missile"
        ],

        "Heavy Assault Missiles": [
            "Mjolnir Heavy Assault Missile",
            "Scourge Heavy Assault Missile",
            "Inferno Heavy Assault Missile",
            "Nova Heavy Assault Missile",

            "Nova Javelin Heavy Assault Missile",
            "Scourge Javelin Heavy Assault Missile",
            "Mjolnir Javelin Heavy Assault Missile",
            "Inferno Javelin Heavy Assault Missile",

            "Nova Rage Heavy Assault Missile",
            "Mjolnir Rage Heavy Assault Missile",
            "Inferno Rage Heavy Assault Missile",
            "Scourge Rage Heavy Assault Missile",

            "Guristas Nova Heavy Assault Missile",
            "Dread Guristas Nova Heavy Assault Missile",
            "Caldari Navy Nova Heavy Assault Missile",
            "Guristas Inferno Heavy Assault Missile",
            "Caldari Navy Inferno Heavy Assault Missile",
            "Dread Guristas Inferno Heavy Assault Missile",
            "Guristas Scourge Heavy Assault Missile",
            "Caldari Navy Scourge Heavy Assault Missile",
            "Dread Guristas Scourge Heavy Assault Missile",
            "Guristas Mjolnir Heavy Assault Missile",
            "Caldari Navy Mjolnir Heavy Assault Missile",
            "Dread Guristas Mjolnir Heavy Assault Missile"
        ],

        "XL Cruise Missiles": [
            "Inferno XL Cruise Missile",
            "Mjolnir XL Cruise Missile",
            "Scourge XL Cruise Missile",
            "Nova XL Cruise Missile",

            "Guristas Scourge XL Cruise Missile",
            "Guristas Nova XL Cruise Missile",
            "Guristas Mjolnir XL Cruise Missile",
            "Guristas Inferno XL Cruise Missile",

            "Scourge Precision XL Cruise Missile",
            "Mjolnir Precision XL Cruise Missile",
            "Inferno Precision XL Cruise Missile",
            "Nova Precision XL Cruise Missile",

            "Scourge Fury XL Cruise Missile",
            "Mjolnir Fury XL Cruise Missile",
            "Nova Fury XL Cruise Missile",
            "Inferno Fury XL Cruise Missile"
        ],

        "Structure Antisubcapital Missiles": [
            "Standup Light Missile",
            "Standup Cruise Missile",
            "Standup Heavy Missile"
        ],

        "Structure Anticapital Missiles": [
            "Standup Super-heavy Torpedo",
            "Standup XL Cruise Missile",
            "Standup AXL-SR Missile"
        ]
    },

    "Cap Booster Charges": [
        "Cap Booster 100",
        "Cap Booster 150",
        "Cap Booster 200",
        "Cap Booster 25",
        "Cap Booster 3200",
        "Cap Booster 400",
        "Cap Booster 50",
        "Cap Booster 75",
        "Cap Booster 800",
        "Navy Cap Booster 100",
        "Navy Cap Booster 150",
        "Navy Cap Booster 200",
        "Navy Cap Booster 25",
        "Navy Cap Booster 3200",
        "Navy Cap Booster 400",
        "Navy Cap Booster 50",
        "Navy Cap Booster 75",
        "Navy Cap Booster 800"
    ],

    "Command Burst Charges": {
        "Structure Anticapital Missiles": [
            "Armor Energizing Charge",
            "Armor Reinforcement Charge",
            "Rapid Repair Charge"
        ],
        "Information Command Burst Charges": [
            "Electronic Hardening Charge",
            "Electronic Superiority Charge",
            "Sensor Optimization Charge"
        ],
        "Mining Foreman Burst Charges": [
            "Mining Equipment Preservation Charge",
            "Mining Laser Field Enhancement Charge",
            "Mining Laser Optimization Charge"

        ],
        "Shield Command Burst Charges": [
            "Active Shielding Charge",
            "Shield Extension Charge",
            "Shield Harmonizing Charge"
        ],
        "Skirmish Command Burst Charges": [
            "Evasive Maneuvers Charge",
            "Interdiction Maneuvers Charge",
            "Rapid Deployment Charge"
        ]
    },

    "Exotic Plasma Charges": {
        "Advanced Exotic Plasma Charges": [
            "Mystic L",
            "Occult L",

            "Mystic M",
            "Occult M",

            "Mystic S",
            "Occult S"
        ],

        "Standard Exotic Plasma Charges": [
            "Baryon Exotic Plasma XL",
            "Meson Exotic Plasma XL",
            "Tetryon Exotic Plasma XL",

            "Baryon Exotic Plasma L",
            "Meson Exotic Plasma L",
            "Tetryon Exotic Plasma L",

            "Baryon Exotic Plasma M",
            "Meson Exotic Plasma M",
            "Tetryon Exotic Plasma M",

            "Baryon Exotic Plasma S",
            "Meson Exotic Plasma S",
            "Tetryon Exotic Plasma S"
        ]
    },

    "Frequency Crystals": {
        "Advanced Beam Laser Crystals": [
            "Aurora XL",
            "Gleam XL",

            "Aurora L",
            "Gleam L",

            "Aurora M",
            "Gleam M",

            "Aurora S",
            "Gleam S"
        ],

        "Advanced Pulse Laser Crystals": [
            "Conflagration XL",
            "Scorch XL",

            "Conflagration L",
            "Scorch L",

            "Conflagration M",
            "Scorch M",

            "Conflagration S",
            "Scorch S"
        ],

        "Faction Crystals": [
            "Blood Gamma XL",
            "Blood Infrared XL",
            "Blood Microwave XL",
            "Blood Multifrequency XL",
            "Blood Radio XL",
            "Blood Standard XL",
            "Blood Ultraviolet XL",
            "Blood Xray XL",
            "Dark Blood Gamma XL",
            "Dark Blood Infrared XL",
            "Dark Blood Microwave XL",
            "Dark Blood Multifrequency XL",
            "Dark Blood Radio XL",
            "Dark Blood Standard XL",
            "Dark Blood Ultraviolet XL",
            "Dark Blood Xray XL",
            "Sanshas Gamma XL",
            "Sanshas Infrared XL",
            "Sanshas Microwave XL",
            "Sanshas Multifrequency XL",
            "Sanshas Radio XL",
            "Sanshas Standard XL",
            "Sanshas Ultraviolet XL",
            "Sanshas Xray XL",
            "True Sanshas Gamma XL",
            "True Sanshas Infrared XL",
            "True Sanshas Microwave XL",
            "True Sanshas Multifrequency XL",
            "True Sanshas Radio XL",
            "True Sanshas Standard XL",
            "True Sanshas Ultraviolet XL",
            "True Sanshas Xray XL",

            "Blood Gamma L",
            "Blood Infrared L",
            "Blood Microwave L",
            "Blood Multifrequency L",
            "Blood Radio L",
            "Blood Standard L",
            "Blood Ultraviolet L",
            "Blood Xray L",
            "Dark Blood Gamma L",
            "Dark Blood Infrared L",
            "Dark Blood Microwave L",
            "Dark Blood Multifrequency L",
            "Dark Blood Radio L",
            "Dark Blood Standard L",
            "Dark Blood Ultraviolet L",
            "Dark Blood Xray L",
            "Imperial Navy Gamma L",
            "Imperial Navy Infrared L",
            "Imperial Navy Microwave L",
            "Imperial Navy Multifrequency L",
            "Imperial Navy Radio L",
            "Imperial Navy Standard L",
            "Imperial Navy Ultraviolet L",
            "Imperial Navy Xray L",
            "Sanshas Gamma L",
            "Sanshas Infrared L",
            "Sanshas Microwave L",
            "Sanshas Multifrequency L",
            "Sanshas Radio L",
            "Sanshas Standard L",
            "Sanshas Ultraviolet L",
            "Sanshas Xray L",
            "True Sanshas Gamma L",
            "True Sanshas Infrared L",
            "True Sanshas Microwave L",
            "True Sanshas Multifrequency L",
            "True Sanshas Radio L",
            "True Sanshas Standard L",
            "True Sanshas Ultraviolet L",
            "True Sanshas Xray L",

            "Blood Gamma M",
            "Blood Infrared M",
            "Blood Microwave M",
            "Blood Multifrequency M",
            "Blood Radio M",
            "Blood Standard M",
            "Blood Ultraviolet M",
            "Blood Xray M",
            "Dark Blood Gamma M",
            "Dark Blood Infrared M",
            "Dark Blood Microwave M",
            "Dark Blood Multifrequency M",
            "Dark Blood Radio M",
            "Dark Blood Standard M",
            "Dark Blood Ultraviolet M",
            "Dark Blood Xray M",
            "Imperial Navy Gamma M",
            "Imperial Navy Infrared M",
            "Imperial Navy Microwave M",
            "Imperial Navy Multifrequency M",
            "Imperial Navy Radio M",
            "Imperial Navy Standard M",
            "Imperial Navy Ultraviolet M",
            "Imperial Navy Xray M",
            "Sanshas Gamma M",
            "Sanshas Infrared M",
            "Sanshas Microwave M",
            "Sanshas Multifrequency M",
            "Sanshas Radio M",
            "Sanshas Standard M",
            "Sanshas Ultraviolet M",
            "Sanshas Xray M",
            "True Sanshas Gamma M",
            "True Sanshas Infrared M",
            "True Sanshas Microwave M",
            "True Sanshas Multifrequency M",
            "True Sanshas Radio M",
            "True Sanshas Standard M",
            "True Sanshas Ultraviolet M",
            "True Sanshas Xray M",

            "Blood Gamma S",
            "Blood Infrared S",
            "Blood Microwave S",
            "Blood Multifrequency S",
            "Blood Radio S",
            "Blood Standard S",
            "Blood Ultraviolet S",
            "Blood Xray S",
            "Dark Blood Gamma S",
            "Dark Blood Infrared S",
            "Dark Blood Microwave S",
            "Dark Blood Multifrequency S",
            "Dark Blood Radio S",
            "Dark Blood Standard S",
            "Dark Blood Ultraviolet S",
            "Dark Blood Xray S",
            "Imperial Navy Gamma S",
            "Imperial Navy Infrared S",
            "Imperial Navy Microwave S",
            "Imperial Navy Multifrequency S",
            "Imperial Navy Radio S",
            "Imperial Navy Standard S",
            "Imperial Navy Ultraviolet S",
            "Imperial Navy Xray S",
            "Sanshas Gamma S",
            "Sanshas Infrared S",
            "Sanshas Microwave S",
            "Sanshas Multifrequency S",
            "Sanshas Radio S",
            "Sanshas Standard S",
            "Sanshas Ultraviolet S",
            "Sanshas Xray S",
            "True Sanshas Gamma S",
            "True Sanshas Infrared S",
            "True Sanshas Microwave S",
            "True Sanshas Multifrequency S",
            "True Sanshas Radio S",
            "True Sanshas Standard S",
            "True Sanshas Ultraviolet S",
            "True Sanshas Xray S"
        ],

        "Standard Crystals": [
            "Gamma XL",
            "Infrared XL",
            "Microwave XL",
            "Multifrequency XL",
            "Radio XL",
            "Standard XL",
            "Ultraviolet XL",
            "Xray XL",

            "Gamma L",
            "Infrared L",
            "Microwave L",
            "Multifrequency L",
            "Radio L",
            "Standard L",
            "Ultraviolet L",
            "Xray L",

            "Gamma M",
            "Infrared M",
            "Microwave M",
            "Multifrequency M",
            "Radio M",
            "Standard M",
            "Ultraviolet M",
            "Xray M",

            "Gamma S",
            "Infrared S",
            "Microwave S",
            "Multifrequency S",
            "Radio S",
            "Standard S",
            "Ultraviolet S",
            "Xray S"
        ],

        "Orbital Strike": []
    },

    "Hybrid Charges": {
        "Advanced Blaster Charges": [
            "Null XL",
            "Void XL",

            "Null L",
            "Void L",

            "Null M",
            "Void M",

            "Null S",
            "Void S"
        ],

        "Advanced Railgun Charges": [
            "Javelin XL",
            "Spike XL",

            "Javelin L",
            "Spike L",

            "Javelin M",
            "Spike M",

            "Javelin S",
            "Spike S"
        ],

        "Faction Charges": [
            "Dread Guristas Antimatter Charge XL",
            "Dread Guristas Iridium Charge XL",
            "Dread Guristas Iron Charge XL",
            "Dread Guristas Lead Charge XL",
            "Dread Guristas Plutonium Charge XL",
            "Dread Guristas Thorium Charge XL",
            "Dread Guristas Tungsten Charge XL",
            "Dread Guristas Uranium Charge XL",
            "Guardian Antimatter Charge XL",
            "Guardian Iridium Charge XL",
            "Guardian Iron Charge XL",
            "Guardian Lead Charge XL",
            "Guardian Plutonium Charge XL",
            "Guardian Thorium Charge XL",
            "Guardian Tungsten Charge XL",
            "Guardian Uranium Charge XL",
            "Guristas Antimatter Charge XL",
            "Guristas Iridium Charge XL",
            "Guristas Iron Charge XL",
            "Guristas Lead Charge XL",
            "Guristas Plutonium Charge XL",
            "Guristas Thorium Charge XL",
            "Guristas Tungsten Charge XL",
            "Guristas Uranium Charge XL",
            "Shadow Antimatter Charge XL",
            "Shadow Iridium Charge XL",
            "Shadow Iron Charge XL",
            "Shadow Lead Charge XL",
            "Shadow Plutonium Charge XL",
            "Shadow Thorium Charge XL",
            "Shadow Tungsten Charge XL",
            "Shadow Uranium Charge XL",

            "Caldari Navy Antimatter Charge L",
            "Caldari Navy Iridium Charge L",
            "Caldari Navy Iron Charge L",
            "Caldari Navy Lead Charge L",
            "Caldari Navy Plutonium Charge L",
            "Caldari Navy Thorium Charge L",
            "Caldari Navy Tungsten Charge L",
            "Caldari Navy Uranium Charge L",
            "Dread Guristas Antimatter Charge L",
            "Dread Guristas Iridium Charge L",
            "Dread Guristas Iron Charge L",
            "Dread Guristas Lead Charge L",
            "Dread Guristas Plutonium Charge L",
            "Dread Guristas Thorium Charge L",
            "Dread Guristas Tungsten Charge L",
            "Dread Guristas Uranium Charge L",
            "Federation Navy Antimatter Charge L",
            "Federation Navy Iridium Charge L",
            "Federation Navy Iron Charge L",
            "Federation Navy Lead Charge L",
            "Federation Navy Plutonium Charge L",
            "Federation Navy Thorium Charge L",
            "Federation Navy Tungsten Charge L",
            "Federation Navy Uranium Charge L",
            "Guardian Antimatter Charge L",
            "Guardian Iridium Charge L",
            "Guardian Iron Charge L",
            "Guardian Lead Charge L",
            "Guardian Plutonium Charge L",
            "Guardian Thorium Charge L",
            "Guardian Tungsten Charge L",
            "Guardian Uranium Charge L",
            "Guristas Antimatter Charge L",
            "Guristas Iridium Charge L",
            "Guristas Iron Charge L",
            "Guristas Lead Charge L",
            "Guristas Plutonium Charge L",
            "Guristas Thorium Charge L",
            "Guristas Tungsten Charge L",
            "Guristas Uranium Charge L",
            "Shadow Antimatter Charge L",
            "Shadow Iridium Charge L",
            "Shadow Iron Charge L",
            "Shadow Lead Charge L",
            "Shadow Plutonium Charge L",
            "Shadow Thorium Charge L",
            "Shadow Tungsten Charge L",
            "Shadow Uranium Charge L",

            "Caldari Navy Antimatter Charge M",
            "Caldari Navy Iridium Charge M",
            "Caldari Navy Iron Charge M",
            "Caldari Navy Lead Charge M",
            "Caldari Navy Plutonium Charge M",
            "Caldari Navy Thorium Charge M",
            "Caldari Navy Tungsten Charge M",
            "Caldari Navy Uranium Charge M",
            "Dread Guristas Antimatter Charge M",
            "Dread Guristas Iridium Charge M",
            "Dread Guristas Iron Charge M",
            "Dread Guristas Lead Charge M",
            "Dread Guristas Plutonium Charge M",
            "Dread Guristas Thorium Charge M",
            "Dread Guristas Tungsten Charge M",
            "Dread Guristas Uranium Charge M",
            "Federation Navy Antimatter Charge M",
            "Federation Navy Iridium Charge M",
            "Federation Navy Iron Charge M",
            "Federation Navy Lead Charge M",
            "Federation Navy Plutonium Charge M",
            "Federation Navy Thorium Charge M",
            "Federation Navy Tungsten Charge M",
            "Federation Navy Uranium Charge M",
            "Guardian Antimatter Charge M",
            "Guardian Iridium Charge M",
            "Guardian Iron Charge M",
            "Guardian Lead Charge M",
            "Guardian Plutonium Charge M",
            "Guardian Thorium Charge M",
            "Guardian Tungsten Charge M",
            "Guardian Uranium Charge M",
            "Guristas Antimatter Charge M",
            "Guristas Iridium Charge M",
            "Guristas Iron Charge M",
            "Guristas Lead Charge M",
            "Guristas Plutonium Charge M",
            "Guristas Thorium Charge M",
            "Guristas Tungsten Charge M",
            "Guristas Uranium Charge M",
            "Shadow Antimatter Charge M",
            "Shadow Iridium Charge M",
            "Shadow Iron Charge M",
            "Shadow Lead Charge M",
            "Shadow Plutonium Charge M",
            "Shadow Thorium Charge M",
            "Shadow Tungsten Charge M",
            "Shadow Uranium Charge M",

            "Caldari Navy Antimatter Charge S",
            "Caldari Navy Iridium Charge S",
            "Caldari Navy Iron Charge S",
            "Caldari Navy Lead Charge S",
            "Caldari Navy Plutonium Charge S",
            "Caldari Navy Thorium Charge S",
            "Caldari Navy Tungsten Charge S",
            "Caldari Navy Uranium Charge S",
            "Dread Guristas Antimatter Charge S",
            "Dread Guristas Iridium Charge S",
            "Dread Guristas Iron Charge S",
            "Dread Guristas Lead Charge S",
            "Dread Guristas Plutonium Charge S",
            "Dread Guristas Thorium Charge S",
            "Dread Guristas Tungsten Charge S",
            "Dread Guristas Uranium Charge S",
            "Federation Navy Antimatter Charge S",
            "Federation Navy Iridium Charge S",
            "Federation Navy Iron Charge S",
            "Federation Navy Lead Charge S",
            "Federation Navy Plutonium Charge S",
            "Federation Navy Thorium Charge S",
            "Federation Navy Tungsten Charge S",
            "Federation Navy Uranium Charge S",
            "Guardian Antimatter Charge S",
            "Guardian Iridium Charge S",
            "Guardian Iron Charge S",
            "Guardian Lead Charge S",
            "Guardian Plutonium Charge S",
            "Guardian Thorium Charge S",
            "Guardian Tungsten Charge S",
            "Guardian Uranium Charge S",
            "Guristas Antimatter Charge S",
            "Guristas Iridium Charge S",
            "Guristas Iron Charge S",
            "Guristas Lead Charge S",
            "Guristas Plutonium Charge S",
            "Guristas Thorium Charge S",
            "Guristas Tungsten Charge S",
            "Guristas Uranium Charge S",
            "Shadow Antimatter Charge S",
            "Shadow Iridium Charge S",
            "Shadow Iron Charge S",
            "Shadow Lead Charge S",
            "Shadow Plutonium Charge S",
            "Shadow Thorium Charge S",
            "Shadow Tungsten Charge S",
            "Shadow Uranium Charge S"
        ],

        "Standard Charges": [
            "Antimatter Charge XL",
            "Iridium Charge XL",
            "Iron Charge XL",
            "Lead Charge XL",
            "Plutonium Charge XL",
            "Thorium Charge XL",
            "Tungsten Charge XL",
            "Uranium Charge XL",

            "Antimatter Charge L",
            "Iridium Charge L",
            "Iron Charge L",
            "Lead Charge L",
            "Plutonium Charge L",
            "Thorium Charge L",
            "Tungsten Charge L",
            "Uranium Charge L",

            "Antimatter Charge M",
            "Iridium Charge M",
            "Iron Charge M",
            "Lead Charge M",
            "Plutonium Charge M",
            "Thorium Charge M",
            "Tungsten Charge M",
            "Uranium Charge M",

            "Antimatter Charge S",
            "Iridium Charge S",
            "Iron Charge S",
            "Lead Charge S",
            "Plutonium Charge S",
            "Thorium Charge S",
            "Tungsten Charge S",
            "Uranium Charge S"
        ],

        "Orbital Strike": []
    },

    "Mining Crystals": {
        "Asteroid Mining Crystals": [
            "Abyssal Asteroid Mining Crystal Type A I",
            "Abyssal Asteroid Mining Crystal Type A II",
            "Abyssal Asteroid Mining Crystal Type B I",
            "Abyssal Asteroid Mining Crystal Type B II",
            "Abyssal Asteroid Mining Crystal Type C I",
            "Abyssal Asteroid Mining Crystal Type C II",
            "Coherent Asteroid Mining Crystal Type A I",
            "Coherent Asteroid Mining Crystal Type A II",
            "Coherent Asteroid Mining Crystal Type B I",
            "Coherent Asteroid Mining Crystal Type B II",
            "Coherent Asteroid Mining Crystal Type C I",
            "Coherent Asteroid Mining Crystal Type C II",
            "Complex Asteroid Mining Crystal Type A I",
            "Complex Asteroid Mining Crystal Type A II",
            "Complex Asteroid Mining Crystal Type B I",
            "Complex Asteroid Mining Crystal Type B II",
            "Complex Asteroid Mining Crystal Type C I",
            "Complex Asteroid Mining Crystal Type C II",
            "Mercoxit Asteroid Mining Crystal Type A I",
            "Mercoxit Asteroid Mining Crystal Type A II",
            "Mercoxit Asteroid Mining Crystal Type B I",
            "Mercoxit Asteroid Mining Crystal Type B II",
            "Mercoxit Asteroid Mining Crystal Type C I",
            "Mercoxit Asteroid Mining Crystal Type C II",
            "Simple Asteroid Mining Crystal Type A I",
            "Simple Asteroid Mining Crystal Type A II",
            "Simple Asteroid Mining Crystal Type B I",
            "Simple Asteroid Mining Crystal Type B II",
            "Simple Asteroid Mining Crystal Type C I",
            "Simple Asteroid Mining Crystal Type C II",
            "Variegated Asteroid Mining Crystal Type A I",
            "Variegated Asteroid Mining Crystal Type A II",
            "Variegated Asteroid Mining Crystal Type B I",
            "Variegated Asteroid Mining Crystal Type B II",
            "Variegated Asteroid Mining Crystal Type C I",
            "Variegated Asteroid Mining Crystal Type C II"

        ],
        "Moon Mining Crystals": [
            "Common Moon Mining Crystal Type A I",
            "Common Moon Mining Crystal Type A II",
            "Common Moon Mining Crystal Type B I",
            "Common Moon Mining Crystal Type B II",
            "Common Moon Mining Crystal Type C I",
            "Common Moon Mining Crystal Type C II",
            "Exceptional Moon Mining Crystal Type A I",
            "Exceptional Moon Mining Crystal Type A II",
            "Exceptional Moon Mining Crystal Type B I",
            "Exceptional Moon Mining Crystal Type B II",
            "Exceptional Moon Mining Crystal Type C I",
            "Exceptional Moon Mining Crystal Type C I",
            "Rare Moon Mining Crystal Type A I",
            "Rare Moon Mining Crystal Type A II",
            "Rare Moon Mining Crystal Type B I",
            "Rare Moon Mining Crystal Type B II",
            "Rare Moon Mining Crystal Type C I",
            "Rare Moon Mining Crystal Type C II",
            "Ubiquitous Moon Mining Crystal Type A I",
            "Ubiquitous Moon Mining Crystal Type A II",
            "Ubiquitous Moon Mining Crystal Type B I",
            "Ubiquitous Moon Mining Crystal Type B II",
            "Ubiquitous Moon Mining Crystal Type C I",
            "Ubiquitous Moon Mining Crystal Type C II",
            "Uncommon Moon Mining Crystal Type A I",
            "Uncommon Moon Mining Crystal Type A II",
            "Uncommon Moon Mining Crystal Type B I",
            "Uncommon Moon Mining Crystal Type B II",
            "Uncommon Moon Mining Crystal Type C I",
            "Uncommon Moon Mining Crystal Type C II"
        ]
    },

    "Nanite Repair Paste": ["Nanite Repair Paste"],

    "Probes": {
        "Interdiction Probes": [
            "Stasis Webification Probe",
            "Surgical Warp Disrupt Probe",
            "Warp Disrupt Probe"
        ],

        "Scan Probes": [
            "Combat Scanner Probe I",
            "Core Scanner Probe I",
            "RSS Core Scanner Probe",
            "Satori-Horigu Combat Scanner Probe",
            "Sisters Combat Scanner Probe",
            "Sisters Core Scanner Probe"
        ],

        "Survey Probes": [
            "Discovery Survey Probe I",
            "Gaze Survey Probe I",
            "Quest Survey Probe I"
        ]
    },

    "Projectile Ammo": {
        "Advanced Artillery Ammo": [
            "Quake XL",
            "Tremor XL",

            "Quake L",
            "Tremor L",

            "Quake M",
            "Tremor M",

            "Quake S",
            "Tremor S"
        ],

        "Advanced Autocannon Ammo": [
            "Barrage XL",
            "Hail XL",

            "Barrage L",
            "Hail L",

            "Barrage M",
            "Hail M",

            "Barrage S",
            "Hail S"
        ],

        "Faction Ammo": [
            "Arch Angel Carbonized Lead XL",
            "Arch Angel Depleted Uranium XL",
            "Arch Angel EMP XL",
            "Arch Angel Fusion XL",
            "Arch Angel Nuclear XL",
            "Arch Angel Phased Plasma XL",
            "Arch Angel Proton XL",
            "Arch Angel Titanium Sabot XL",
            "Domination Carbonized Lead XL",
            "Domination Depleted Uranium XL",
            "Domination EMP XL",
            "Domination Fusion XL",
            "Domination Nuclear XL",
            "Domination Phased Plasma XL",
            "Domination Proton XL",
            "Domination Titanium Sabot XL",
            "Republic Fleet Carbonized Lead XL",
            "Republic Fleet Depleted Uranium XL",
            "Republic Fleet EMP XL",
            "Republic Fleet Fusion XL",
            "Republic Fleet Nuclear XL",
            "Republic Fleet Phased Plasma XL",
            "Republic Fleet Proton XL",
            "Republic Fleet Titanium Sabot XL",

            "Arch Angel Carbonized Lead L",
            "Arch Angel Depleted Uranium L",
            "Arch Angel EMP L",
            "Arch Angel Fusion L",
            "Arch Angel Nuclear L",
            "Arch Angel Phased Plasma L",
            "Arch Angel Proton L",
            "Arch Angel Titanium Sabot L",
            "Domination Carbonized Lead L",
            "Domination Depleted Uranium L",
            "Domination EMP L",
            "Domination Fusion L",
            "Domination Nuclear L",
            "Domination Phased Plasma L",
            "Domination Proton L",
            "Domination Titanium Sabot L",
            "Republic Fleet Carbonized Lead L",
            "Republic Fleet Depleted Uranium L",
            "Republic Fleet EMP L",
            "Republic Fleet Fusion L",
            "Republic Fleet Nuclear L",
            "Republic Fleet Phased Plasma L",
            "Republic Fleet Proton L",
            "Republic Fleet Titanium Sabot L",

            "Arch Angel Carbonized Lead M",
            "Arch Angel Depleted Uranium M",
            "Arch Angel EMP M",
            "Arch Angel Fusion M",
            "Arch Angel Nuclear M",
            "Arch Angel Phased Plasma M",
            "Arch Angel Proton M",
            "Arch Angel Titanium Sabot M",
            "Domination Carbonized Lead M",
            "Domination Depleted Uranium M",
            "Domination EMP M",
            "Domination Fusion M",
            "Domination Nuclear M",
            "Domination Phased Plasma M",
            "Domination Proton M",
            "Domination Titanium Sabot M",
            "Republic Fleet Carbonized Lead M",
            "Republic Fleet Depleted Uranium M",
            "Republic Fleet EMP M",
            "Republic Fleet Fusion M",
            "Republic Fleet Nuclear M",
            "Republic Fleet Phased Plasma M",
            "Republic Fleet Proton M",
            "Republic Fleet Titanium Sabot M",

            "Arch Angel Carbonized Lead S",
            "Arch Angel Depleted Uranium S",
            "Arch Angel EMP S",
            "Arch Angel Fusion S",
            "Arch Angel Nuclear S",
            "Arch Angel Phased Plasma S",
            "Arch Angel Proton S",
            "Arch Angel Titanium Sabot S",
            "Domination Carbonized Lead S",
            "Domination Depleted Uranium S",
            "Domination EMP S",
            "Domination Fusion S",
            "Domination Nuclear S",
            "Domination Phased Plasma S",
            "Domination Proton S",
            "Domination Titanium Sabot S",
            "Republic Fleet Carbonized Lead S",
            "Republic Fleet Depleted Uranium S",
            "Republic Fleet EMP S",
            "Republic Fleet Fusion S",
            "Republic Fleet Nuclear S",
            "Republic Fleet Phased Plasma S",
            "Republic Fleet Proton S",
            "Republic Fleet Titanium Sabot S"
        ],

        "Standard Ammo": [
            "Carbonized Lead XL",
            "Depleted Uranium XL",
            "EMP XL",
            "Fusion XL",
            "Nuclear XL",
            "Phased Plasma XL",
            "Proton XL",
            "Titanium Sabot XL",

            "Carbonized Lead L",
            "Depleted Uranium L",
            "EMP L",
            "Fusion L",
            "Nuclear L",
            "Phased Plasma L",
            "Proton L",
            "Titanium Sabot L",

            "Carbonized Lead M",
            "Depleted Uranium M",
            "EMP M",
            "Fusion M",
            "Nuclear M",
            "Phased Plasma M",
            "Proton M",
            "Titanium Sabot M",

            "Carbonized Lead S",
            "Depleted Uranium S",
            "EMP S",
            "Fusion S",
            "Nuclear S",
            "Phased Plasma S",
            "Proton S",
            "Titanium Sabot S"
        ],

        "Orbital Strike": []
    },

    "Scripts": [
        "Armor EM Resistance Script",
        "Armor Explosive Resistance Script",
        "Armor Kinetic Resistance Script",
        "Armor Thermal Resistance Script",
        "ECCM Script",
        "Focused Warp Disruption Script",
        "Focused Warp Scrambling Script",
        "Missile Precision Disruption Script",
        "Missile Precision Script",
        "Missile Range Disruption Script",
        "Missile Range Script",
        "Optimal Range Disruption Script",
        "Optimal Range Script",
        "Scan Resolution Dampening Script",
        "Scan Resolution Script",
        "Shield EM Resistance Script",
        "Shield Explosive Resistance Script",
        "Shield Kinetic Resistance Script",
        "Shield Thermal Resistance Script",
        "Standup Focused Warp Scrambling Script",
        "Standup Gravimetric ECM Script",
        "Standup Ladar ECM Script",
        "Standup Magnetometric ECM Script",
        "Standup Radar ECM Script",
        "Targeting Range Dampening Script",
        "Targeting Range Script",
        "Tracking Speed Disruption Script",
        "Tracking Speed Script"
    ],

    "Structure Area Denial Ammunition": ["Standup Flak Round I"],
    "Structure Guided Bombs": ["Standup Heavy Guided Bomb", "Standup Light Guided Bomb"]
};

const emit = defineEmits(['selectItem']);
 
const selectItem = (item: string) => {
    emit('selectItem', item);
}

const toggleExpanded = (path: string) => {
    const newSet = new Set(expandedPaths.value);
    if (newSet.has(path)) {
        newSet.delete(path);
    } else {
        newSet.add(path);
    }
    expandedPaths.value = newSet;
};
</script>