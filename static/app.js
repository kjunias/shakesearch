const Model = {
  batch: 0,
  results: [],
}

const Controller = {
  search: (ev) => {
    ev.preventDefault();
    Model.results = [];
    Model.batch = 0;
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    Model.results = Model.results.concat(results);
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of Model.results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    table.innerHTML = rows;
    Model.batch++;
  },

  loadMore: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&b=${Model.batch}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  }
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

const button = document.getElementById("load-more");
button.addEventListener("click", Controller.loadMore);
