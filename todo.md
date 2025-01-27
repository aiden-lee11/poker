basics is that we need a five card evaluator function 


we also need a simulation function that uses the 5 card eval with simulated turns and rivers


for pre flop we should just use precomputed values and not simulate because its a known quantity and i dont wanna be simulating 5 random cards

as for once flop is down we just have to simulate remaining two, iterate through the 7 chose 5 combinations of a players cards pick the highest strength out of them and then finally compare to highest strength of the other players 

each simulation step keep track of who won and then at the end of all simulations just do the basic win percentage for equity i guess


when flop is the only thing down we need to simulate turn and river

when flop and turn is down only simulate river

when all common cards are down no simulation is needed just evaluate the players 5 best cards out of the available 7 and return the highest value as the winning player


links http://suffe.cool/poker/evaluator.html

https://github.com/christophschmalhofer/poker/blob/master/XPokerEval/XPokerEval.CactusKev.PerfectHash/fast_eval.cpp


