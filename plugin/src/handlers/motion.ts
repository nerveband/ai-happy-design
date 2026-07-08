import { getSceneNodeById } from '../utils/getNode';

function unavailable(feature: string): never {
  throw new Error(feature + ' is unavailable in this Figma runtime');
}

export async function handleMotion(action: string, params: any): Promise<any> {
  switch (action) {
    case 'get_styles': return getStyles();
    case 'apply_style': return applyStyle(params);
    case 'remove_style': return removeStyle(params);
    case 'get_animations': return getAnimations(params);
    case 'apply_keyframes': return applyKeyframes(params);
    case 'remove_keyframes': return removeKeyframes(params);
    case 'set_timeline_duration': return setTimelineDuration(params);
    default: throw new Error('Unknown motion action: ' + action + '. Available: get_styles, apply_style, remove_style, get_animations, apply_keyframes, remove_keyframes, set_timeline_duration');
  }
}

async function getStyles() {
  var figmaAny = figma as any;
  if (typeof figmaAny.getLocalMotionStylesAsync !== 'function') unavailable('Motion styles');
  var styles = await figmaAny.getLocalMotionStylesAsync();
  return { styles: JSON.parse(JSON.stringify(styles || [])), count: styles ? styles.length : 0 };
}

async function applyStyle(params: any) {
  var node = await getSceneNodeById(params.nodeId);
  var anyNode = node as any;
  if (typeof anyNode.setMotionStyleIdAsync !== 'function') unavailable('Motion style application');
  await anyNode.setMotionStyleIdAsync(params.styleId || '');
  return { id: node.id, name: node.name, motionStyleId: anyNode.motionStyleId || params.styleId || '' };
}

async function removeStyle(params: any) {
  var node = await getSceneNodeById(params.nodeId);
  var anyNode = node as any;
  if (typeof anyNode.setMotionStyleIdAsync !== 'function') unavailable('Motion style removal');
  await anyNode.setMotionStyleIdAsync('');
  return { id: node.id, name: node.name, motionStyleId: '' };
}

async function getAnimations(params: any) {
  var node = await getSceneNodeById(params.nodeId);
  var anyNode = node as any;
  if (!('animations' in anyNode) && !('motionKeyframes' in anyNode)) unavailable('Motion animations');
  return {
    id: node.id,
    name: node.name,
    animations: JSON.parse(JSON.stringify(anyNode.animations || anyNode.motionKeyframes || [])),
    timelineDuration: anyNode.timelineDuration,
  };
}

async function applyKeyframes(params: any) {
  var node = await getSceneNodeById(params.nodeId);
  var anyNode = node as any;
  if (!('motionKeyframes' in anyNode) && typeof anyNode.setMotionKeyframes !== 'function') unavailable('Motion keyframes');
  var keyframes = typeof params.keyframes === 'string' ? JSON.parse(params.keyframes) : params.keyframes;
  if (!Array.isArray(keyframes)) throw new Error('keyframes must be an array');
  if ('motionKeyframes' in anyNode) anyNode.motionKeyframes = keyframes;
  else anyNode.setMotionKeyframes(keyframes);
  if (params.duration != null && 'timelineDuration' in anyNode) anyNode.timelineDuration = params.duration;
  return { id: node.id, name: node.name, keyframeCount: keyframes.length, timelineDuration: anyNode.timelineDuration };
}

async function removeKeyframes(params: any) {
  var node = await getSceneNodeById(params.nodeId);
  var anyNode = node as any;
  if (!('motionKeyframes' in anyNode) && typeof anyNode.setMotionKeyframes !== 'function') unavailable('Motion keyframes');
  if ('motionKeyframes' in anyNode) anyNode.motionKeyframes = [];
  else anyNode.setMotionKeyframes([]);
  return { id: node.id, name: node.name, keyframeCount: 0 };
}

async function setTimelineDuration(params: any) {
  var figmaAny = figma as any;
  if (!('motionTimelineDuration' in figmaAny) && !('timelineDuration' in figmaAny)) unavailable('Motion timeline duration');
  if ('motionTimelineDuration' in figmaAny) figmaAny.motionTimelineDuration = params.duration;
  else figmaAny.timelineDuration = params.duration;
  return { timelineDuration: params.duration };
}
