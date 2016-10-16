// Graph Labels
var names = ['Queries', 'Abnormal'];
var groups = new vis.DataSet();

//Initialize Groups
groups.add({
  id: 0,
  content: names[0],
  options: {
    drawPoints: false,
    interpolation: false
  }
});

groups.add({
  id: 1,
  content: names[1],
  options: {
    drawPoints: false,
    interpolation: false
  }
});

// Chart id in Body
var container = document.getElementById('visualization');
var dataset = new vis.DataSet();

//Option For Group
var options = {
  dataAxis: {
    left: {
      range: {
        min: 0
      },
    },
    showMinorLabels: false
  },
  drawPoints: false,
  legend: true,
  start: vis.moment().add(-30, 'seconds'),
  end: vis.moment(),
};


var graph2d = new vis.Graph2d(container, dataset, groups, options);

//Initilase Queries & Abnormal Values
var yValuesOld = {
  Total: 0,
  Abnormal: 0
}
var yValues = {
  Total: 0,
  Abnormal: 0
}

//Return diff of last Two Queries
function yQueries() {
  console.log(yValues["Total"] - yValuesOld["Total"])
  return yValues["Total"] - yValuesOld["Total"];
}

//Return diff of last Two Abnormals
function yAbnormal() {
  console.log("X")
  console.log(yValues["Abnormal"] - yValuesOld["Abnormal"])
  return yValues["Abnormal"] - yValuesOld["Abnormal"];
}

function renderStep() {
  // move the window (you can think of different strategies).
  var now = vis.moment();
  var range = graph2d.getWindow();
  var interval = range.end - range.start;

  // move the window 90% to the left when now is larger than the end of the window
  if (now > range.end) {
    //graph2d.setWindow(now - interval, now, {animation: true});
    graph2d.setWindow(now - 0.1 * interval, now + 0.9 * interval, {
      animation: true
    });
  }
  setTimeout(renderStep, 1000);

}


//Add datapoint to the graph
function addDataPoint() {
  var now = vis.moment();
  dataset.add([{
    x: now,
    y: yQueries(),
    group: 0
  }, {
    x: now,
    y: yAbnormal(),
    group: 1
  }]);
  setTimeout(addDataPoint, 1000);
}

//Get Value from Api
function loadDoc() {

  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      yValuesOld = yValues;
      yValues = JSON.parse(this.responseText);
    }
  };
  xhttp.open("GET", '/api?_=' + new Date().getTime(), true);
  xhttp.send();
}

renderStep();
addDataPoint();
loadDoc();
setInterval(loadDoc, 1000);
