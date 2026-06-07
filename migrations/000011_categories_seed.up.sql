-- Очищаємо старі дані (зберігаємо платежі, просто обнуляємо категорію)
UPDATE payments SET category_id = NULL;
DELETE FROM payment_categories;
DELETE FROM category_groups;

-- ── Групи витрат ──────────────────────────────────────────────────
INSERT INTO category_groups (name, is_income) VALUES
  ('Їжа та напої',        false),
  ('Покупки',             false),
  ('Житло',               false),
  ('Транспорт',           false),
  ('Автотранспорт',       false),
  ('Життя та розваги',    false),
  ('Зв''язок',            false),
  ('Фінансові витрати',   false),
  ('Інвестиції',          false),
  ('Інше',                false),
-- ── Групи доходів ─────────────────────────────────────────────────
  ('Дохід',               true);

-- ── Категорії витрат ──────────────────────────────────────────────
INSERT INTO payment_categories (name, icon, group_id) VALUES
  -- Їжа та напої
  ('Їжа та напої',                 '🍱', (SELECT id FROM category_groups WHERE name = 'Їжа та напої')),
  ('Бар, кафе',                    '☕', (SELECT id FROM category_groups WHERE name = 'Їжа та напої')),
  ('Ресторан, фаст-фуд',           '🍔', (SELECT id FROM category_groups WHERE name = 'Їжа та напої')),
  ('Продукти',                     '🛒', (SELECT id FROM category_groups WHERE name = 'Їжа та напої')),

  -- Покупки
  ('Покупки',                      '🛍', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Дрогері та аптека',            '💊', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Вільний час',                  '🎯', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Канцелярія, інструменти',      '📎', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Подарунки',                    '🎁', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Радощі',                       '😊', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Електроніка',                  '💻', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Домашні тварини',              '🐾', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Дім, сад',                     '🏡', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Діти',                         '👶', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Здоров''я і краса',            '💅', (SELECT id FROM category_groups WHERE name = 'Покупки')),
  ('Прикраси, аксесуари',          '💍', (SELECT id FROM category_groups WHERE name = 'Покупки')),

  -- Житло
  ('Страхування майна',            '🔒', (SELECT id FROM category_groups WHERE name = 'Житло')),
  ('Житло',                        '🏠', (SELECT id FROM category_groups WHERE name = 'Житло')),
  ('Обслуговування',               '🔧', (SELECT id FROM category_groups WHERE name = 'Житло')),
  ('Ремонт',                       '🪚', (SELECT id FROM category_groups WHERE name = 'Житло')),
  ('Послуги',                      '📋', (SELECT id FROM category_groups WHERE name = 'Житло')),
  ('Комунальні послуги',           '💡', (SELECT id FROM category_groups WHERE name = 'Житло')),
  ('Іпотека',                      '🏦', (SELECT id FROM category_groups WHERE name = 'Житло')),
  ('Оренда',                       '🔑', (SELECT id FROM category_groups WHERE name = 'Житло')),

  -- Транспорт
  ('Транспорт',                    '🚌', (SELECT id FROM category_groups WHERE name = 'Транспорт')),
  ('Відрядження',                  '💼', (SELECT id FROM category_groups WHERE name = 'Транспорт')),
  ('Далекі поїздки',               '✈', (SELECT id FROM category_groups WHERE name = 'Транспорт')),
  ('Таксі',                        '🚕', (SELECT id FROM category_groups WHERE name = 'Транспорт')),
  ('Громадський транспорт',        '🚇', (SELECT id FROM category_groups WHERE name = 'Транспорт')),

  -- Автотранспорт
  ('Страхування транспорту',       '🛡', (SELECT id FROM category_groups WHERE name = 'Автотранспорт')),
  ('Оренда авто',                  '🚗', (SELECT id FROM category_groups WHERE name = 'Автотранспорт')),
  ('Технічне обслуговування',      '🔩', (SELECT id FROM category_groups WHERE name = 'Автотранспорт')),
  ('Паркування',                   '🅿', (SELECT id FROM category_groups WHERE name = 'Автотранспорт')),
  ('Паливо',                       '⛽', (SELECT id FROM category_groups WHERE name = 'Автотранспорт')),

  -- Життя та розваги
  ('Життя та розваги',             '🎭', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Лотереї, азартні ігри',        '🎰', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Алкоголь, тютюн',              '🍷', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Благодійність, подарунки',     '❤',  (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Відпустка, подорожі, готелі',  '🏖', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Освіта, розвиток',             '📚', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Хобі',                         '🎨', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Життєві події',                '🎉', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Культура, спортивні події',    '🎪', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Активний спорт, фітнес',       '🏃', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Велнес, краса',                '🧖', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),
  ('Медична допомога, лікар',      '🏥', (SELECT id FROM category_groups WHERE name = 'Життя та розваги')),

  -- Зв'язок
  ('Зв''язок, ПК',                 '💬', (SELECT id FROM category_groups WHERE name = 'Зв''язок')),
  ('Поштові послуги',              '📮', (SELECT id FROM category_groups WHERE name = 'Зв''язок')),
  ('Програми, застосунки, ігри',   '📱', (SELECT id FROM category_groups WHERE name = 'Зв''язок')),
  ('Інтернет',                     '🌐', (SELECT id FROM category_groups WHERE name = 'Зв''язок')),
  ('Телефонія, мобільний',         '📞', (SELECT id FROM category_groups WHERE name = 'Зв''язок')),

  -- Фінансові витрати
  ('Фінансові витрати',            '💸', (SELECT id FROM category_groups WHERE name = 'Фінансові витрати')),
  ('Аліменти',                     '👨‍👩‍👧', (SELECT id FROM category_groups WHERE name = 'Фінансові витрати')),
  ('Комісії, збори',               '🏷', (SELECT id FROM category_groups WHERE name = 'Фінансові витрати')),
  ('Консультації',                 '🤝', (SELECT id FROM category_groups WHERE name = 'Фінансові витрати')),
  ('Штрафи',                       '⚠',  (SELECT id FROM category_groups WHERE name = 'Фінансові витрати')),
  ('Кредити, відсотки',            '📉', (SELECT id FROM category_groups WHERE name = 'Фінансові витрати')),
  ('Податки',                      '🧾', (SELECT id FROM category_groups WHERE name = 'Фінансові витрати')),

  -- Інвестиції
  ('Інвестиції',                   '📈', (SELECT id FROM category_groups WHERE name = 'Інвестиції')),
  ('Колекції',                     '🏺', (SELECT id FROM category_groups WHERE name = 'Інвестиції')),
  ('Заощадження',                  '🐖', (SELECT id FROM category_groups WHERE name = 'Інвестиції')),
  ('Фінансові інвестиції',         '💹', (SELECT id FROM category_groups WHERE name = 'Інвестиції')),
  ('Транспорт, рухоме майно',      '🚢', (SELECT id FROM category_groups WHERE name = 'Інвестиції')),

  -- Інше
  ('Відсутньо',                    '❓', (SELECT id FROM category_groups WHERE name = 'Інше')),
  ('Інше',                         '📦', (SELECT id FROM category_groups WHERE name = 'Інше'));

-- ── Категорії доходів ─────────────────────────────────────────────
INSERT INTO payment_categories (name, icon, group_id) VALUES
  ('Дохід',                        '💰', (SELECT id FROM category_groups WHERE name = 'Дохід')),
  ('Подарунки',                    '🎁', (SELECT id FROM category_groups WHERE name = 'Дохід')),
  ('Аліменти',                     '👨‍👩‍👧', (SELECT id FROM category_groups WHERE name = 'Дохід')),
  ('Лотереї, азартні ігри',        '🎰', (SELECT id FROM category_groups WHERE name = 'Дохід')),
  ('Дохід від позики',             '💳', (SELECT id FROM category_groups WHERE name = 'Дохід')),
  ('Внески і гранти',              '🤲', (SELECT id FROM category_groups WHERE name = 'Дохід')),
  ('Дохід від оренди',             '🏘', (SELECT id FROM category_groups WHERE name = 'Дохід')),
  ('Продаж',                       '🛒', (SELECT id FROM category_groups WHERE name = 'Дохід')),
  ('Відсотки, дивіденди',          '📊', (SELECT id FROM category_groups WHERE name = 'Дохід')),
  ('Зарплата',                     '💵', (SELECT id FROM category_groups WHERE name = 'Дохід'));

-- Фіналізуємо NOT NULL
ALTER TABLE payment_categories ALTER COLUMN group_id SET NOT NULL;
