const adjectives = [
  'happy', 'cosmic', 'neon', 'fuzzy', 'sparkly', 'groovy', 'stellar', 'electric',
  'golden', 'mystic', 'radiant', 'vivid', 'bold', 'swift', 'bright', 'jolly',
  'dazzle', 'turbo', 'mega', 'super', 'ultra', 'hyper', 'epic', 'lucky',
  'wild', 'brave', 'noble', 'pixel', 'retro', 'cyber', 'astral', 'coral',
];

const nouns = [
  'unicorn', 'rainbow', 'rocket', 'pizza', 'penguin', 'dragon', 'taco', 'octopus',
  'phoenix', 'panda', 'koala', 'owl', 'falcon', 'dolphin', 'tiger', 'wolf',
  'narwhal', 'waffle', 'comet', 'nebula', 'quasar', 'prism', 'aurora', 'blaze',
  'meteor', 'orbit', 'spark', 'thunder', 'vortex', 'zenith', 'frost', 'cactus',
];

export function generateChannelKey(): string {
  const adj = adjectives[Math.floor(Math.random() * adjectives.length)];
  const noun = nouns[Math.floor(Math.random() * nouns.length)];
  const num = Math.floor(Math.random() * 100);
  return `${adj}-${noun}-${num}`;
}
