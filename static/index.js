const fetchItems = async () => {
  const response = await fetch('/items');
  const items = await response.json();

  const ul = document.createElement('ul');

  items.forEach((item) => {
    const li = document.createElement('li')
    li.textContent = item.Name;
    li.dataset.uuid = item.UUID;
    if (item.Done) {
      li.className = 'done';
    }
    ul.append(li);
  });

  document.querySelector('#items').append(ul);
}

const addNewItem = async (evt) => {
  if (evt.key !== "Enter") return;

  const newItemName = document.querySelector('#newItemName').value;

  const body = {
    Name: newItemName,
  };

  await fetch('/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body),
  });

  window.location.reload();
}

document.addEventListener("DOMContentLoaded", () => {
  fetchItems();

  document.querySelector('#newItemName').addEventListener('keypress', addNewItem);
});
