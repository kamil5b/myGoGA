package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Chromosome struct {
	/*
		GENES :
		[sign 1/-1 ][0-9][0-9][0-9][0-9][sign 1/-1 ][0-9][0-9][0-9][0-9]
		[              X               ][               Y              ]
		SIGN : 0 negative, 1 positive
		MAX : 9.999
		MIN : 9.999
		Fitness : Function(x,y)
	*/
	genes          [10]int
	fitnessFormula func(float32, float32) float32
	x, y, fitness  float32
}

type Generation struct {
	bestChromosome         Chromosome
	population             []Chromosome
	matingProcess          int
	mutationChance         float32
	Xmax, Xmin, Ymax, Ymin float32
}

func (c *Chromosome) fitnessing() {
	//Fitnessing is by calling the function
	c.fitness = c.fitnessFormula(c.x, c.y)
}

func (c *Chromosome) encode() {
	/*
			Encode = x or y to genes
			exmp x : -1.034
		   	check if x < 0 the first index will be -1 and change the x to unsigned
			if x >= 0 the first index will be  1
		    -1.034 -> index 0 = -1 -> -1.034 = 1.034
			1.034 * 1000 = 1034
		    1034 % 10 = 4 -> index 4 = 4-> 1034 / 10 = 103
		    103 % 10 = 3 -> index 3 = 3 -> 103 / 10 = 10
			10 % 10 = 0 -> index 2 = 0 -> 10/10 = 1
		    1 % 10 = 1 -> index 1 = 1.
			therefore, the half chromosome will be : [-1][1][0][3][4]
			same goes for y
	*/

	//ENCODING X
	temp := int(c.x * 1000)
	if temp < 0 {
		c.genes[0] = -1
		temp *= -1
	} else {
		c.genes[0] = 1
	}

	for i := 0; i < 4; i++ {
		c.genes[4-i] = temp % 10
		temp /= 10
	}
	//ENCODING Y
	temp = int(c.y * 1000)
	if temp < 0 {
		c.genes[5] = -1
		temp *= -1
	} else {
		c.genes[5] = 1
	}

	for i := 0; i < 4; i++ {
		c.genes[9-i] = temp % 10
		temp /= 10
	}
}

func (c *Chromosome) decode() {
	/*
		Decoding : from genes to x y
		simple decode :
		[-1][1][0][3][4] = -1.034
		0.00[4]
		0.0[3]0
		0.[0]00
		[1].000
		--------+
		1.034
		-> 1.034 * [-1] = -1.034
	*/

	x := float32(0)
	y := float32(0)

	for i := 0; i < 5; i++ {
		x += float32(c.genes[i]) * float32(math.Pow(10, float64(-i)))
		y += float32(c.genes[i+5]) * float32(math.Pow(10, float64(-i)))
	}
	x *= float32(c.genes[0])
	y *= float32(c.genes[5])
	c.x = x
	c.y = y
}

func (c *Chromosome) mutate(mutationChance float32) {

	mutation := rand.Float32() * 100
	if mutationChance >= mutation {
		rng := rand.Intn(10)
		idx := rand.Intn(10)
		if idx == 0 || idx == 5 {
			if c.genes[idx] == 1 {
				c.genes[idx] = -1
			} else {
				c.genes[idx] = 1
			}
		} else {
			c.genes[idx] = rng
		}
	}
}

func createChromosome(x, y float32, f func(float32, float32) float32) Chromosome {
	var c Chromosome
	c.x = x
	c.y = y
	c.fitnessFormula = f
	c.fitnessing()
	c.encode()

	return c
}

func createGeneration(totalPopulation, matingProcess int, mutationChance, Xmax, Xmin, Ymax, Ymin float32, f func(float32, float32) float32) Generation {
	var gen Generation

	gen.matingProcess = matingProcess
	gen.mutationChance = mutationChance
	gen.Xmax, gen.Xmin, gen.Ymax, gen.Ymin = Xmax, Xmin, Ymax, Ymin
	population := make([]Chromosome, totalPopulation)

	gen.population = population

	for i := 0; i < totalPopulation; i++ {
		//CREATE RANDOM X and Y with range Xmin to Xmax and Ymin to Ymax
		x := Xmin + rand.Float32()*(Xmax-Xmin)
		y := Ymin + rand.Float32()*(Ymax-Ymin)
		//CREATE THE CHROMOSOME
		gen.population[i] = createChromosome(x, y, f)
	}

	sortChromosomes(gen.population)
	gen.bestChromosome = gen.population[0]

	return gen

}

func sortChromosomes(arrChromosome []Chromosome) {
	//INSERTION SORT
	n := len(arrChromosome)
	for i := 1; i < n; i++ {
		for j := i; j > 0; j-- {
			//DESCENDING, COMPARE FITNESS
			if arrChromosome[j-1].fitness < arrChromosome[j].fitness {
				arrChromosome[j-1], arrChromosome[j] = arrChromosome[j], arrChromosome[j-1]
			}
		}
	}
}

func selectParent(parents []Chromosome) Chromosome {
	//ROULETTE WHEEL

	c := parents[0]
	stop := false
	rouletteWheel := make([]float32, len(parents))
	var pick float32
	totalFitness := float32(0)

	//register first chromosome to the parent list
	if parents[0].fitness < 0 {
		rouletteWheel[0] = parents[0].fitness * -1
	} else {
		rouletteWheel[0] = parents[0].fitness
	}
	totalFitness += parents[0].fitness

	//Register the rest of chromosome to the parent list
	for i := 1; i < len(parents); i++ {
		if parents[i].fitness < 0 {
			rouletteWheel[i] = rouletteWheel[i-1] + parents[i].fitness*-1
		} else {
			rouletteWheel[i] = rouletteWheel[i-1] + parents[i].fitness
		}
		totalFitness += parents[i].fitness
	}

	//PICK RANDOM FLOAT NUMBER FROM 0 to MAX FITNESS
	for i := 0; i < 10; i++ {
		pick = rand.Float32() * totalFitness
	}

	//if the picked number inside the range of the fitness value, than pick the chromosome
	temp := float32(0)
	for i := 0; i < len(parents) && !stop; i++ {
		temp = temp + rouletteWheel[i]
		if pick <= temp {
			c = parents[i]
			stop = true
		}
	}
	return c

}

func matingChromosome(parent1, parent2 Chromosome, Xmax, Xmin, Ymax, Ymin, mutationChance float32) Chromosome {
	//EACH MATING RESULTING THE BEST CHROMOSOME THEIR PRODUCED

	var bestChild, child Chromosome
	//THERE WILL BE 16 SETS
	//each sets are described in the spreadsheet that i made
	//each mating will had 16 child and we pick the best set (bruteforcing)
	child.fitnessFormula = parent1.fitnessFormula
	//SET 1
	for i := 0; i < 5; i++ {
		//C X = P1 X
		child.genes[i] = parent1.genes[i]
		//C Y = P2 Y
		child.genes[i+5] = parent2.genes[i+5]
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		bestChild = child
	}

	//SET 2
	for i := 0; i < 5; i++ {
		//C X = P2 X
		child.genes[i] = parent2.genes[i]
		//C Y = P1 Y
		child.genes[i+5] = parent1.genes[i+5]
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 3
	for i := 0; i < 5; i++ {
		//C X = P2 Y
		child.genes[i] = parent2.genes[i+5]
		//C Y = P1 X
		child.genes[i+5] = parent1.genes[i]
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 4
	for i := 0; i < 5; i++ {
		//C X = P1 Y
		child.genes[i] = parent1.genes[i+5]
		//C Y = P2 X
		child.genes[i+5] = parent2.genes[i]
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 5
	for i := 0; i < 5; i++ {
		//C X = P1 (odd) P2 (even) X
		if i%2 != 0 {
			//ODD P1
			child.genes[i] = parent1.genes[i]
		} else {
			//EVEN P2
			child.genes[i] = parent2.genes[i]
		}
		//C Y = P2 Y
		child.genes[i+5] = parent2.genes[i+5]
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 6
	for i := 0; i < 5; i++ {
		//C X = P2 X
		child.genes[i] = parent2.genes[i]
		//C Y = P1 (ODD) P2 (EVEN) Y
		if i%2 != 0 {
			//ODD P1
			child.genes[i+5] = parent1.genes[i+5]
		} else {
			//EVEN P2
			child.genes[i+5] = parent2.genes[i+5]
		}
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 7
	for i := 0; i < 5; i++ {
		//C X = P2 (odd) P1 (even) X
		if i%2 != 0 {
			//ODD P2
			child.genes[i] = parent2.genes[i]
		} else {
			//EVEN P1
			child.genes[i] = parent1.genes[i]
		}
		//C Y = P2 Y
		child.genes[i+5] = parent2.genes[i+5]
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 8
	for i := 0; i < 5; i++ {
		//C X = P2 X
		child.genes[i] = parent2.genes[i]
		//C Y = P2 (ODD) P1 (EVEN) Y
		if i%2 != 0 {
			//ODD P2
			child.genes[i+5] = parent2.genes[i+5]
		} else {
			//EVEN P1
			child.genes[i+5] = parent1.genes[i+5]
		}
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 9
	for i := 0; i < 5; i++ {
		//C X = P1 (odd) P2 (even) X
		if i%2 != 0 {
			//ODD P1
			child.genes[i] = parent1.genes[i]
		} else {
			//EVEN P2
			child.genes[i] = parent2.genes[i]
		}
		//C Y = P1 Y
		child.genes[i+5] = parent1.genes[i+5]
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 10
	for i := 0; i < 5; i++ {
		//C X = P1 X
		child.genes[i] = parent1.genes[i]
		//C Y = P1 (ODD) P2 (EVEN) Y
		if i%2 != 0 {
			//ODD P1
			child.genes[i+5] = parent1.genes[i+5]
		} else {
			//EVEN P2
			child.genes[i+5] = parent2.genes[i+5]
		}
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 11
	for i := 0; i < 5; i++ {
		//C X = P2 (odd) P1 (even) X
		if i%2 != 0 {
			//ODD P2
			child.genes[i] = parent2.genes[i]
		} else {
			//EVEN P1
			child.genes[i] = parent1.genes[i]
		}
		//C Y = P1 Y
		child.genes[i+5] = parent1.genes[i+5]
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 12
	for i := 0; i < 5; i++ {
		//C X = P1 X
		child.genes[i] = parent1.genes[i]
		//C Y = P2 (ODD) P1 (EVEN) Y
		if i%2 != 0 {
			//ODD P2
			child.genes[i+5] = parent2.genes[i+5]
		} else {
			//EVEN P1
			child.genes[i+5] = parent1.genes[i+5]
		}
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 13
	for i := 0; i < 5; i++ {
		//C X = P1 (odd) P2 (even) X
		if i%2 != 0 {
			//ODD P1
			child.genes[i] = parent1.genes[i]
		} else {
			//EVEN P2
			child.genes[i] = parent2.genes[i]
		}
		//C Y = P1 (ODD) P2 (EVEN) Y
		if i%2 != 0 {
			//ODD P1
			child.genes[i+5] = parent1.genes[i+5]
		} else {
			//EVEN P2
			child.genes[i+5] = parent2.genes[i+5]
		}
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 14
	for i := 0; i < 5; i++ {
		//C X = P2 (odd) P1 (even) X
		if i%2 != 0 {
			//ODD P2
			child.genes[i] = parent2.genes[i]
		} else {
			//EVEN P1
			child.genes[i] = parent1.genes[i]
		}
		//C Y = P2 (ODD) P1 (EVEN) Y
		if i%2 != 0 {
			//ODD P2
			child.genes[i+5] = parent2.genes[i+5]
		} else {
			//EVEN P1
			child.genes[i+5] = parent1.genes[i+5]
		}
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 15
	for i := 0; i < 5; i++ {
		//C X = P1 (odd) P2 (even) X
		if i%2 != 0 {
			//ODD P1
			child.genes[i] = parent1.genes[i]
		} else {
			//EVEN P2
			child.genes[i] = parent2.genes[i]
		}
		//C Y = P2 (ODD) P1 (EVEN) Y
		if i%2 != 0 {
			//ODD P2
			child.genes[i+5] = parent2.genes[i+5]
		} else {
			//EVEN P1
			child.genes[i+5] = parent1.genes[i+5]
		}
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	//SET 16
	for i := 0; i < 5; i++ {
		//C X = P2 (odd) P1 (even) X
		if i%2 != 0 {
			//ODD P2
			child.genes[i] = parent2.genes[i]
		} else {
			//EVEN P1
			child.genes[i] = parent1.genes[i]
		}
		//C Y = P1 (ODD) P2 (EVEN) Y
		if i%2 != 0 {
			//ODD P1
			child.genes[i+5] = parent1.genes[i+5]
		} else {
			//EVEN P2
			child.genes[i+5] = parent2.genes[i+5]
		}
	}
	child.mutate(mutationChance)
	child.decode()
	if child.x >= Xmin && child.x <= Xmax && child.y >= Ymin && child.y <= Ymax {
		child.fitnessing()
		if bestChild.fitness <= child.fitness {
			bestChild = child
		}
	}

	return bestChild
}

func (g Generation) reGeneration() {
	matingPool := g.population
	for i := 0; i < g.matingProcess; i++ {
		parent1 := selectParent(matingPool)
		parent2 := selectParent(matingPool)

		child := matingChromosome(parent1, parent2, g.Xmax, g.Xmin, g.Ymax, g.Ymin, g.mutationChance)

		matingPool = append(matingPool, child)

	}
	sortChromosomes(matingPool)
	totalPopulation := len(g.population)
	g.population = matingPool[:totalPopulation]
	g.bestChromosome = matingPool[0]
}

func generateGenerations(generations, totalPopulation, matingProcess int, mutationChance, Xmax, Xmin, Ymax, Ymin float32, f func(float32, float32) float32) Generation {
	gen := createGeneration(totalPopulation, matingProcess, mutationChance, Xmax, Xmin, Ymax, Ymin, f)
	for i := 0; i < generations; i++ {
		gen.reGeneration()
	}
	return gen
}

func (c Chromosome) viewChromosome() {
	for i := 0; i < 10; i++ {
		fmt.Print("[", c.genes[i], "]")
	}
	fmt.Println()
	fmt.Println("x : ", c.x)
	fmt.Println("y : ", c.y)
	fmt.Println("fitness value : ", c.fitness)
}

func (g Generation) viewGeneration() {
	for i := 0; i < len(g.population); i++ {
		fmt.Println("------", i+1, "------")
		g.population[i].viewChromosome()
		fmt.Println("------------------")
	}
	fmt.Println("Best Chromosome :")
	g.bestChromosome.viewChromosome()
}
