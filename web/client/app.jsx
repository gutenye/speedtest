var {LineChart, BarChart} = ReactD3
_ = lodash
pd = console.log.bind(console)

var PROVIDERS = ["vultr-tokyo", "vultr-seattle", "vultr-sydney", "vultr-silicon", "vultr-la", "vultr-amsterdam", "vultr-frankfurt", "vultr-london", "vultr-paris", "vultr-chicago", "vultr-dallas", "vultr-nyc", "vultr-atlanta", "vultr-miami", "do-nyc1", "do-nyc2", "do-nyc3", "do-sgp1", "do-ams1", "do-ams2"]

var App = React.createClass({
  mixins: [ ReactiveMixin ],

  getReactiveState: function(arg) {
    var data = []
    arg = arg || {}
    var playload = Result.find(arg).fetch()
    if (playload.length === 0) {
      return {data: []}
    }
    playload.forEach(v => {
      var found = _.find(data, "label", v.name)
      var item = null
      if (found) {
        item = found
      } else {
        item = {label: v.name, values: []}
        data.push(item)
      }
      item.values.push({x: v.date, y: v.avg / 1000})
    })

    return {
      data: data,
      xScale: d3.time.scale().domain([playload[0].date, playload[playload.length-1].date]).range([0, 800])
    }
  },

  render() {
    if (this.state.data.length === 0 ) {
      return <div />
    }
    var {data, xScale} = this.state
    var tooltipHtml = function(label, data) {
      //return `${label} x: ${data.x} y: ${data.y}`
      return `${label} ${data.y}`
    }

    return (
      <div>
        <LineChart data={data} width={870} height={400} margin={{top: 10, bottom: 50, left: 50, right: 20}} tooltipHtml={tooltipHtml} xScale={xScale} />
        {PROVIDERS.map(v => {
          return <div onClick={this.change.bind(this, v)}>{v}</div>
        })}
      </div>
    )
  },

  change(name) {
    var state = this.getReactiveState({name: name})
    this.setState(state)
  }
})

Meteor.startup(function() {
  React.render(<App />, document.body)
})
