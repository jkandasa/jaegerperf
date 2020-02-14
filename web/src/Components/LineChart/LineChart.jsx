import React from "react"
import { Line } from "react-chartjs-2"
import "chartjs-plugin-colorschemes"
// https://nagix.github.io/chartjs-plugin-colorschemes/colorchart.html

const lineChart = ({ title = "", data = [], yAxesLabel = "" }) => {
  const optionsChart2 = {
    title: {
      display: false,
      fontSize: 16,
      text: title
    },
    plugins: {
      colorschemes: {
        scheme: "tableau.Classic10"
      }
    },
    maintainAspectRatio: false,
    responsive: true,
    scales: {
      xAxes: [
        {
          offset: true,
          type: "category",
          ticks: {
            source: "data",
            autoSkip: true,
            maxRotation: 45,
            minRotation: 0
          },
          scaleLabel: {
            display: false,
            labelString: "x-label"
          }
        }
      ],
      yAxes: [
        {
          ticks: {
            beginAtZero: true
          },
          scaleLabel: {
            display: true,
            labelString: yAxesLabel
          }
        }
      ]
    },
    tooltips: {
      mode: "index",
      intersect: false
    },
    hover: {
      mode: "index",
      intersect: false,
      axis: "xy"
    },
    legend: {
      position: "top"
    },
    elements: {
      point: {
        radius: 0,
        hitRadius: 5,
        hoverRadius: 5
      }
    }
  }

  return (
    <div style={{ minHeight: "370px" }}>
      <Line data={data} options={optionsChart2} />
    </div>
  )
}

export default lineChart
