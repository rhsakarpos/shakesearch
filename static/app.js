const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    const data = Object.fromEntries(new FormData(form));
    for (let result of results.Results) {
      //result = result.replaceAll(data.query, '<b>'+ data.query + '</b>')
      var div = document.createElement('div');
      div.innerHTML = `<h2>${result.Title}</h2>`;
      document.body.appendChild(div);
      var tbl = document.createElement('table');
      tbl.style.width = '100%';
      tbl.setAttribute('border', '1');
      var tbdy = document.createElement('tbody');
      // add rows to the table
      for (let row of result.Matches){
        var tr = document.createElement('tr');
        var td = document.createElement('td');
        //td.appendChild(document.createTextNode(row));
        var span = document.createElement("span");
        tr.appendChild(td);
        td.innerHTML = row;
        tbl.appendChild(tr);
      }
      document.body.appendChild(tbl);
    }

  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
