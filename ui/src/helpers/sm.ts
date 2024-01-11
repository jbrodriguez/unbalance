// interface StateActions {
//   onEnter: () => void;
//   onExit: () => void;
// }

interface Transition {
  target: string;
  action: () => void;
}

interface Transitions {
  [transition: string]: Transition;
}

// interface State {
//   // actions: StateActions;
//   transitions: StateTransitions;
// }

export interface StateMachineDefinition {
  initialState: string;
  [state: string]: Transitions | string;
}

export interface StateMachine {
  value: string;
  transition: (currentState: string, event: string) => string;
}

export const createMachine = (
  stateMachineDefinition: StateMachineDefinition,
): StateMachine => {
  let machineState: string = stateMachineDefinition.initialState;

  const machine: StateMachine = Object.freeze({
    get value() {
      return machineState;
    },
    transition(currentState: string, event: string): string {
      console.log('machine.transition ', currentState, event);
      const currentStateDefinition = stateMachineDefinition[
        currentState
      ] as Transitions;
      const destinationTransition = currentStateDefinition[event];
      if (!destinationTransition) {
        return machineState;
      }
      const destinationState = destinationTransition.target;

      destinationTransition.action();

      machineState = destinationState;

      return machineState;
    },
  });

  return machine;
};

// export const createMachine = (
//   stateMachineDefinition: StateMachineDefinition,
// ): StateMachine => {
//   // eslint-disable-next-line
//   let machine: StateMachine = {
//     value: stateMachineDefinition.initialState,
//     transition(currentState: string, event: string): string {
//       console.log('transition ', currentState, event);
//       const currentStateDefinition = stateMachineDefinition[
//         currentState
//       ] as Transitions;
//       const destinationTransition = currentStateDefinition[event];
//       if (!destinationTransition) {
//         return machine.value;
//       }
//       const destinationState = destinationTransition.target;
//       // const destinationStateDefinition = stateMachineDefinition[
//       //   destinationState
//       // ] as State;

//       destinationTransition.action();
//       // currentStateDefinition.actions.onExit();
//       // destinationStateDefinition.actions.onEnter();

//       machine.value = destinationState;

//       return machine.value;
//     },
//   };
//   return machine as StateMachine;
// };
