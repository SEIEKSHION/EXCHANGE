document.addEventListener('DOMContentLoaded', function() {
	const selectFromValute = document.getElementById('fromValute');
    const selectToValute = document.getElementById('toValute');
    const convertBtn = document.getElementById('convertBtn');
    const resultDiv = document.getElementById('result');
    const errorDiv = document.getElementById('errorMessage');
    
    
    
    // Добавляем класс загрузки
    selectFromValute.classList.add('loading');
    selectToValute.classList.add('loading');
    
    
    // объявление функций:
    function createOption(value, text) {
	    const option = document.createElement('option');
	    option.value = value;
	    option.textContent = text;
	    return option;
	}
    
    
    fetch('/api/valutes')
        .then(response => {
          
            if (!response.ok) {
                throw new Error('Ошибка сети');
            }
            return response.json();
        })
        .then(data => {
            console.log("Полученные данные:", data)
            // удаление класса
            selectFromValute.classList.remove('loading');
            selectToValute.classList.remove('loading');
            
            
            // Добавление элемента - "Выберите валюту"
            selectFromValute.innerHTML = '<option value="">Выберите валюту</option>';
            selectToValute.innerHTML = '<option value="">Выберите валюту</option>';
            
            
            
            
            data.forEach(item => {
			    selectFromValute.appendChild(createOption(item.Value, item.Name));
			    selectToValute.appendChild(createOption(item.Value, item.Name));
			});
			selectFromValute.appendChild(createOption())
			selectToValute.appendChild(createOption())
        })
        .catch(error => {
            console.error('Ошибка:', error);
            
            selectToValute.classList.remove('loading');
            selectFromValute.classList.remove('loading');
            
            errorDiv.textContent = 'Не удалось загрузить список валют. Пожалуйста, попробуйте позже.';
            errorDiv.style.display = 'block';
    });
    convertBtn.addEventListener('click', function() {
    	// взятие данных из формы
        const fromCurrencySelected = selectFromValute.value;
        const toCurrencySelected = selectToValute.value;
        const enteredQuantity = document.getElementById("cQ").valueAsNumber;
        
        // если хоть одна ил валют не выбрана
        if (!fromCurrencySelected || !toCurrencySelected) {
            errorDiv.textContent = 'Пожалуйста, выберите валюту';
            errorDiv.style.display = 'block';
            return;
        }
        
        // Проверка: введено ли количество
        if (enteredQuantity === undefined || enteredQuantity <= 0) {
            errorDiv.textContent = 'Пожалуйста, введите корректное количество';
            errorDiv.style.display = 'block';
            return;
        }
        
        // иначе
        errorDiv.style.display = 'none';
        resultDiv.textContent = `Перевод осуществляется из ${selectFromValute.options[selectFromValute.selectedIndex].text} в ${selectToValute.options[selectToValute.selectedIndex].text}`;
        resultDiv.classList.add('show');
        
        // отпрвляем запрос на сервер
        fetch('/api/convert', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
        	fromcurrency: selectFromValute.options[selectFromValute.selectedIndex].text,
            tocurrency: selectToValute.options[selectToValute.selectedIndex].text,
            quantity: enteredQuantity
            })
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => { throw new Error(err.error || 'Ошибка сервера'); });
            }
            return response.json();
        })
        .then(data => {
            console.log('Результат:', data);

            // ИСПРАВЛЕНО: в тексте был JS-код как строка — теперь корректно
            const fromCurrencyName = selectFromValute.options[selectFromValute.selectedIndex].text;
			const toCurrencyName = selectToValute.options[selectToValute.selectedIndex].text;
			
			resultDiv.textContent = `${enteredQuantity} ${fromCurrencyName} — это ${data.quantity} ${toCurrencyName}`;
            resultDiv.style.display = 'block'; // показываем результат
        })
        .catch(error => {
            console.error('Ошибка:', error);
            errorDiv.textContent = `Ошибка конвертации: ${error.message}`;
            errorDiv.style.display = 'block';
        });
    });
});