import { insertCSS } from '../helpers/insertCSS';
import { hashObject } from '../helpers/hashObject';

/**
 * Animation transition can be used in ReactCSSTransitionGroup
 * */
export interface Transition {
  transitionName: string;
  transitionEnterTimeout: number;
  transitionLeaveTimeout: number;
}

export interface Position {
  x: number;
  y: number;
}

const transitionPrefix = 'tc-15-transition-';
const defaultEnterTimeout = 300;
const defaultLeaveTimeout = 300;

const generateFadeCSS = (transition: Transition, visibleOpacity: number) =>
  `
.${transition.transitionName}-enter {
    opacity: 0 !important;
}
.${transition.transitionName}-enter-active {
    opacity: ${visibleOpacity} !important;
    transition: opacity ${transition.transitionEnterTimeout}ms ease;
}
.${transition.transitionName}-leave {
    opacity: ${visibleOpacity} !important;
}
.${transition.transitionName}-leave-active {
    opacity: 0 !important;
    transition: opacity ${transition.transitionLeaveTimeout}ms ease;
}
`;

const generateSlideCSS = (transition: Transition, enterPosition: Position, leavePosition: Position) =>
  `
.${transition.transitionName}-enter {
    opacity: 0 !important;
    transform: translate3d(${enterPosition.x}px, ${enterPosition.y}px, 0);
}
.${transition.transitionName}-enter-active {
    opacity: 1 !important;
    transform: translate3d(0, 0, 0);
    transition: opacity ${transition.transitionEnterTimeout}ms ease,
                transform ${transition.transitionEnterTimeout}ms ease;
}
.${transition.transitionName}-leave {
    opacity: 1 !important;
    transform: translate3d(0, 0, 0);
}
.${transition.transitionName}-leave-active {
    opacity: 0 !important;
    transform: translate3d(${leavePosition.x}px, ${leavePosition.y}px, 0);
    transition: opacity ${transition.transitionLeaveTimeout}ms ease,
                transform ${transition.transitionLeaveTimeout}ms ease;
}
`;

const generateCollapseCSS = (transition: Transition, targetMaxHeight: number) =>
  `
.${transition.transitionName}-enter {
    max-height: 0;
}
.${transition.transitionName}-enter-active {
    max-height: ${targetMaxHeight}px;
    transition: max-height ${transition.transitionEnterTimeout}ms ease;
}
.${transition.transitionName}-leave {
    max-height: ${targetMaxHeight}px;
}
.${transition.transitionName}-leave-active {
    max-height: 0;
    transition: max-height ${transition.transitionLeaveTimeout}ms ease;
}
`;

let setupTransitions: { [hash: number]: Transition } = {};

/**
 * generate a fade transition, enjoy it with ReactCSSTransitionGroup
 * */
export function fade(visibleOpacity = 1, enterTimeout = defaultEnterTimeout, leaveTimeout = defaultLeaveTimeout) {
  const name = transitionPrefix + 'fade';
  const hash = hashObject({ name, visibleOpacity, enterTimeout, leaveTimeout });

  if (!setupTransitions[hash]) {
    const transition = (setupTransitions[hash] = {
      transitionName: name + hash,
      transitionEnterTimeout: enterTimeout,
      transitionLeaveTimeout: leaveTimeout + 0.0001
    });
    const css = generateFadeCSS(transition, visibleOpacity);
    insertCSS(name + hash, css);
  }

  return setupTransitions[hash];
}

/**
 * generate a slide transition, enjoy it with ReactCSSTransitionGroup
 * */
export function slide(
  enterPosition: Position = { x: 0, y: -30 },
  leavePosition: Position = enterPosition,
  enterTimeout = defaultEnterTimeout,
  leaveTimeout = defaultLeaveTimeout
) {
  const name = transitionPrefix + 'slide';
  const hash = hashObject({ name, enterPosition, leavePosition, enterTimeout, leaveTimeout });

  if (!setupTransitions[hash]) {
    const transition = (setupTransitions[hash] = {
      transitionName: name + hash,
      transitionEnterTimeout: enterTimeout,
      transitionLeaveTimeout: leaveTimeout
    });
    const css = generateSlideCSS(transition, enterPosition, leavePosition);
    insertCSS(name + hash, css);
  }

  return setupTransitions[hash];
}

/**
 * generate a slide transition, enjoy it with ReactCSSTransitionGroup
 * */
export function collapse(targetHeight: number, enterTimeout = defaultEnterTimeout, leaveTimeout = defaultLeaveTimeout) {
  const name = transitionPrefix + 'collapse';
  const hash = hashObject({ name, targetHeight, enterTimeout, leaveTimeout });

  if (!setupTransitions[hash]) {
    const transition = (setupTransitions[hash] = {
      transitionName: name + hash,
      transitionEnterTimeout: enterTimeout,
      transitionLeaveTimeout: leaveTimeout
    });
    const css = generateCollapseCSS(transition, targetHeight);
    insertCSS(name + hash, css);
  }

  return setupTransitions[hash];
}
