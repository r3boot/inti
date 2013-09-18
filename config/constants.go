package config

const MAX_CHANNELS int = 512        // Maximum nr of channels on a DMX bus
const MAX_GROUPS int = 4096         // Maximum nr of groups
const MAX_GROUP_MEMBERS int = 4096  // Maximum nr of members in a group
const MAX_CONTROLLERS int = 4096    // Maximum nr of controllers

const CHAN_FEAT_PWM uint8 = 0x1     // Channel is PWM based
const CHAN_FEAT_ONOFF uint8 = 0x2   // Channel is ON/OFF based
