const fetchItems = async () => {
  const response = await fetch('/items');
  const items = await response.json();

  const ul = document.createElement('ul');

  items.forEach((item) => {
    const li = document.createElement('li');
    const span = document.createElement('span');
    const doneButton = document.createElement('button');
    const removeButton = document.createElement('button');

    span.textContent = item.Name;
    doneButton.textContent = 'Done';
    doneButton.addEventListener('click', toggleDone);
    removeButton.textContent = 'Remove';
    removeButton.addEventListener('click', removeItem);

    li.append(doneButton);
    li.append(removeButton);
    li.append(span);
    li.dataset.uuid = item.UUID;
    li.dataset.done = item.Done;

    if (item.Done) {
      li.className = 'done';
    }

    ul.append(li);
  });

  document.querySelector('#items').append(ul);
};

const addNewItem = async (evt) => {
  if (evt.key !== "Enter") return;

  const newItemName = document.querySelector('#newItemName').value;

  const body = {
    Name: newItemName,
  };

  await fetch('/items/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body),
  });

  window.location.reload();
};

const toggleDone = async (evt) => {
  const { uuid, done } = evt.target.parentElement.dataset;

  const body = {
    Done: done === 'true' ? false : true
  };

  await fetch(`/items/update/${uuid}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body),
  });

  window.location.reload();
};

const removeItem = async (evt) => {
  const { uuid } = evt.target.parentElement.dataset;

  await fetch(`/items/delete/${uuid}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json'
    },
  });

  window.location.reload();
};

document.addEventListener("DOMContentLoaded", () => {
  fetchItems();

  document.querySelector('#newItemName').addEventListener('keypress', addNewItem);
});
