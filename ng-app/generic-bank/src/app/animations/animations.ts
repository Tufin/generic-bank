import {animate, group, query, stagger, state, style, transition, trigger} from "@angular/animations";

const ANIMATION_DURATION = '0.3s';
const ANIMATION_SHORT_DURATION = '0.2s';
const ANIMATION_TYPE = 'ease-in-out';

const navigationOpacity = trigger('navigationOpacity', [
  transition('* <=> *', [
    /* order */
    /* 1 */ query(':enter', style({ opacity: '0'})
      , { optional: true }),
    /* 2 */ group([  // block executes in parallel
      query(':enter', [
        style({ opacity: '0' }),
        animate(`${ANIMATION_DURATION} ${ANIMATION_TYPE}`, style({ opacity: '1'}))
      ], { optional: true }),
      query(':leave', [
        style({ opacity: '1' }),
        animate(`${ANIMATION_DURATION} ${ANIMATION_TYPE}`, style({ opacity: '0' }))
      ], { optional: true }),
    ])
  ])
]);
//
// const opacityTransition = trigger('opacityTransition', [
//   transition('* <=> *', [
//     // query(':self, :enter', [
//       query(':self:enter', [
//         style({ opacity: '0' }),
//         animate(`${ANIMATION_DURATION} ${ANIMATION_TYPE}`, style({ opacity: '1'}))
//       ], { optional: true }),
//       query(':self:leave', [
//         style({ opacity: '1', position: 'absolute', width: 'calc(100% - 30px)' }),
//         animate(`${ANIMATION_DURATION} ${ANIMATION_TYPE}`, style({ opacity: '0' }))
//       ], { optional: true })
//     ])
//   // ])
// ]);

const routerTransition = trigger('routerTransition', [
  transition('* <=> *', [
    /* order */
    /* 1 */ query(':enter, :leave', style({ position: 'absolute'})
      , { optional: true }),
    /* 2 */ group([  // block executes in parallel
      query(':enter', [
        style({ transform: 'translateX(100%)' }),
        animate(`${ANIMATION_DURATION} ${ANIMATION_TYPE}`, style({ transform: 'translateX(0%)' }))
      ], { optional: true }),
      query(':leave', [
        style({ transform: 'translateX(0%)' }),
        animate(`${ANIMATION_DURATION} ${ANIMATION_TYPE}`, style({ transform: 'translateX(-100%)' }))
      ], { optional: true }),
    ])
  ])
]);

const expandTransition = trigger('expandTransition', [
  transition('* <=> *', [
    /* order */
    /* 1 */ query(':enter .card', style({ opacity: 0 })
      , { optional: true }),
    /* 2 */ group([  // block executes in parallel
      query(':enter .card', stagger(150, [
        style({ transform: 'translateY(-50%)' , opacity: 0}),
        animate(`${ANIMATION_SHORT_DURATION} ${ANIMATION_TYPE}`, style({ transform: 'translateY(0%)', opacity: 1 }))
      ]), { optional: true }),
      query(':leave .card', stagger(150, [
        style({ transform: 'translateY(0%)' , opacity: 1}),
        animate(`${ANIMATION_SHORT_DURATION} ${ANIMATION_TYPE}`, style({ transform: 'translateY(50%)', opacity: 0 }))
      ]), { optional: true }),
    ])
  ])
]);

export default {routerTransition, expandTransition, navigationOpacity};
