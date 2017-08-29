
//expects nodes and edges variable

var cy = cytoscape({
  container: document.querySelector('#cy'),

  boxSelectionEnabled: false,
  autounselectify: true,

  style: cytoscape.stylesheet()
    .selector('node')
      .css({
        'content': 'data(name)',
        'text-valign': 'center',
        'color': 'data(color)',
        'text-outline-width': 2,
        'background-color': '#999',
        'text-outline-color': '#999'
      })
    .selector('edge')
      .css({
        'curve-style': 'bezier',
        'target-arrow-shape': 'triangle',
        'target-arrow-color': '#ccc',
        'line-color': '#ccc',
        'width': 1
      })
    .selector(':selected')
      .css({
        'background-color': 'black',
        'line-color': 'black',
        'target-arrow-color': 'black',
        'source-arrow-color': 'black'
      })
    .selector('.faded')
      .css({
        'opacity': 0.25,
        'text-opacity': 0
      }),

  elements: {
    nodes: nodes,
    edges: edges
  },

  layout: {
    name: 'cose',
    idealEdgeLength: 100,
    nodeOverlap: 20
  }
});

cy.on('tap', 'node', function(e){
  var node = e.target;
  var neighborhood = node.neighborhood().add(node);

  cy.elements().addClass('faded');
  neighborhood.removeClass('faded');
});

cy.on('tap', function(e){
  if( e.target === cy ){
    cy.elements().removeClass('faded');
  }
});
