package main

import (
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

const (
	WIDTH = 1000
	HEIGHT = 800
)

var (
	document	js.Value
	canvas		js.Value
	context		js.Value
)

func init() {
	document = js.Global().Get("document")
	canvas = document.Call("getElementById", "mycanvas")
	context = canvas.Call("getContext", "2d")

	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {

	game := new(Game)

	var main_callback js.Callback

	main_callback = js.NewCallback(func(args []js.Value) {
		game.Iterate()
		game.Draw(args)
		js.Global().Call("requestAnimationFrame", main_callback)
	})

	js.Global().Call("requestAnimationFrame", main_callback)

	// We apparently need to stop the Go app from reaching the end,
	// this is one way to do it...

	done := make(chan bool, 0)
	<- done
}

// -------------------------------------------------------------------------
// Graphics helper functions

func set(x, y int) {
	context.Call("fillRect", x, y, 1, 1)
}

func frect(x1, y1, x2, y2 int) {
	if x2 < x1 { x1, x2 = x2, x1 }
	if y2 < y1 { y1, y2 = y2, y1 }
	context.Call("fillRect", x1, y1, x2 - x1, y2 - y1)
}

// -------------------------------------------------------------------------
// Game logic

const (
    QUEENS = 8
    BEASTS = 1900
    BEAST_MAX_SPEED = 7
    QUEEN_MAX_SPEED = 5.5
    BEAST_ACCEL_MODIFIER = 0.55
    QUEEN_ACCEL_MODIFIER = 0.7
    QUEEN_TURN_PROB = 0.001
    BEAST_TURN_PROB = 0.002
    AVOID_STRENGTH = 4000
    MAX_PLAYER_SPEED = 10
    MARGIN = 50
)

const (
    QUEEN = iota
    BEAST
)

type Game struct {
	inited				bool
	queens  			[]*Dood
	beasts				[]*Dood
}

type Dood struct {
	x float64
	y float64
	speedx float64
	speedy float64
	species int
	target *Dood
	game *Game
}

func (d *Dood) Move() {

    x, y, speedx, speedy := d.x, d.y, d.speedx, d.speedy

    var turnprob, maxspeed, accelmod float64
    switch d.species {
    case QUEEN:
        turnprob = QUEEN_TURN_PROB
        maxspeed = QUEEN_MAX_SPEED
        accelmod = QUEEN_ACCEL_MODIFIER
    case BEAST:
        turnprob = BEAST_TURN_PROB
        maxspeed = BEAST_MAX_SPEED
        accelmod = BEAST_ACCEL_MODIFIER
    }

    // Chase target...

    if d.target == nil || rand.Float64() < turnprob || d.target == d {
        tar_id := rand.Intn(QUEENS)
        d.target = d.game.queens[tar_id]
    }

    vecx, vecy := unit_vector(x, y, d.target.x, d.target.y)

    if vecx == 0 && vecy == 0 {
        speedx += (rand.Float64() * 2 - 1) * accelmod
        speedy += (rand.Float64() * 2 - 1) * accelmod
    } else {
        speedx += vecx * rand.Float64() * accelmod
        speedy += vecy * rand.Float64() * accelmod
    }

    // Wall avoidance...

    if (x < MARGIN) {
        speedx += rand.Float64() * 2
    }
    if (x >= WIDTH - MARGIN) {
        speedx -= rand.Float64() * 2
    }
    if (y < MARGIN) {
        speedy += rand.Float64() * 2
    }
    if (y >= HEIGHT - MARGIN) {
        speedy -= rand.Float64() * 2
    }

    // Throttle speed...

    speed := math.Sqrt(speedx * speedx + speedy * speedy)

    if speed > maxspeed {
        speedx *= maxspeed / speed
        speedy *= maxspeed / speed
    }

    // Update entity...

    d.speedx = speedx
    d.speedy = speedy
    d.x += speedx
    d.y += speedy
}

func (self *Game) Draw(args []js.Value) {

	context.Set("fillStyle", "rgb(0,0,0)")
	frect(0, 0, WIDTH, HEIGHT)

	context.Set("fillStyle", "rgb(0,255,0)")
	for _, beast := range self.beasts {
		set(int(beast.x), int(beast.y))
	}
}

func (self *Game) Iterate() {
	if self.inited == false {
		self.Init()
	}

	for _, queen := range self.queens {
		queen.Move()
	}

	for _, beast := range self.beasts {
		beast.Move()
	}
}

func (self *Game) Init() {

	for i := 0 ; i < QUEENS ; i++ {
		self.queens = append(self.queens, &Dood{
											x: rand.Float64() * WIDTH,
											y: rand.Float64() * HEIGHT,
											species: QUEEN,
											game: self,
		})
	}

	for i := 0 ; i < BEASTS ; i++ {
		self.beasts = append(self.beasts, &Dood{
											x: rand.Float64() * WIDTH,
											y: rand.Float64() * HEIGHT,
											species: BEAST,
											game: self,
		})
	}

	self.inited = true
}

// -------------------------------------------------------------------------
// Other helper functions

func unit_vector(x1, y1, x2, y2 float64) (float64, float64) {
    dx := x2 - x1
    dy := y2 - y1

    if (dx < 0.01 && dx > -0.01 && dy < 0.01 && dy > -0.01) {
        return 0, 0
    }

    distance := math.Sqrt(dx * dx + dy * dy)
    return dx / distance, dy / distance
}
