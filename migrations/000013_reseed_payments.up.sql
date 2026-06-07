DELETE FROM payments;

DO $$
DECLARE
  wid INTEGER;

  -- витрати
  c_продукти    INTEGER;
  c_кафе        INTEGER;
  c_ресторан    INTEGER;
  c_комунальні  INTEGER;
  c_таксі       INTEGER;
  c_транспорт   INTEGER;
  c_електроніка INTEGER;
  c_спорт       INTEGER;
  c_освіта      INTEGER;
  c_медицина    INTEGER;
  c_телефон     INTEGER;
  c_ігри        INTEGER;

  -- доходи
  c_зарплата    INTEGER;
  c_відсотки    INTEGER;
  c_продаж      INTEGER;

  -- методи
  m_готівка     INTEGER;
  m_дебетова    INTEGER;
  m_кредитна    INTEGER;
  m_переказ     INTEGER;
  m_мобільний   INTEGER;

BEGIN
  SELECT id INTO wid FROM wallets ORDER BY id DESC LIMIT 1;
  IF wid IS NULL THEN
    RAISE EXCEPTION 'Спочатку зареєструйся — гаманець ще не створено';
  END IF;

  -- категорії витрат
  SELECT id INTO c_продукти    FROM payment_categories WHERE name = 'Продукти';
  SELECT id INTO c_кафе        FROM payment_categories WHERE name = 'Бар, кафе';
  SELECT id INTO c_ресторан    FROM payment_categories WHERE name = 'Ресторан, фаст-фуд';
  SELECT id INTO c_комунальні  FROM payment_categories WHERE name = 'Комунальні послуги';
  SELECT id INTO c_таксі       FROM payment_categories WHERE name = 'Таксі';
  SELECT id INTO c_транспорт   FROM payment_categories WHERE name = 'Громадський транспорт';
  SELECT id INTO c_електроніка FROM payment_categories WHERE name = 'Електроніка';
  SELECT id INTO c_спорт       FROM payment_categories WHERE name = 'Активний спорт, фітнес';
  SELECT id INTO c_освіта      FROM payment_categories WHERE name = 'Освіта, розвиток';
  SELECT id INTO c_медицина    FROM payment_categories WHERE name = 'Медична допомога, лікар';
  SELECT id INTO c_телефон     FROM payment_categories WHERE name = 'Телефонія, мобільний';
  SELECT id INTO c_ігри        FROM payment_categories WHERE name = 'Програми, застосунки, ігри';

  -- категорії доходів
  SELECT id INTO c_зарплата  FROM payment_categories WHERE name = 'Зарплата';
  SELECT id INTO c_відсотки  FROM payment_categories WHERE name = 'Відсотки, дивіденди';
  SELECT id INTO c_продаж    FROM payment_categories WHERE name = 'Продаж';

  -- методи
  SELECT id INTO m_готівка   FROM payment_methods WHERE name = 'Готівка';
  SELECT id INTO m_дебетова  FROM payment_methods WHERE name = 'Дебетова карта';
  SELECT id INTO m_кредитна  FROM payment_methods WHERE name = 'Кредитна карта';
  SELECT id INTO m_переказ   FROM payment_methods WHERE name = 'Банківський переказ';
  SELECT id INTO m_мобільний FROM payment_methods WHERE name = 'Мобільний платіж';

  IF c_зарплата IS NULL OR m_готівка IS NULL THEN
    RAISE EXCEPTION 'Категорії або методи не знайдено — спочатку запусти міграції 000011 та 000012';
  END IF;

  -- ── Поточний місяць ───────────────────────────────────────────
  INSERT INTO payments (wallet_id, is_income, amount, category_id, payment_method_id, description, transaction_date) VALUES
    (wid, true,  48000.00, c_зарплата,   m_переказ,   'Зарплата за місяць',            NOW() - INTERVAL '2 days'),
    (wid, false,  1950.00, c_продукти,   m_дебетова,  'Сільпо — тижневий запас',       NOW() - INTERVAL '1 day'),
    (wid, false,   380.00, c_кафе,       m_мобільний, 'Ранкова кава з колегою',        NOW() - INTERVAL '3 days'),
    (wid, false,   720.00, c_ресторан,   m_кредитна,  'Обід у суші-барі',              NOW() - INTERVAL '4 days'),
    (wid, false,  1400.00, c_комунальні, m_переказ,   'Квартплата + електрика',        NOW() - INTERVAL '5 days'),
    (wid, false,   560.00, c_таксі,      m_мобільний, 'Таксі до аеропорту',           NOW() - INTERVAL '6 days'),
    (wid, false,   140.00, c_транспорт,  m_дебетова,  'Поповнення карти метро',        NOW() - INTERVAL '7 days'),
    (wid, true,   9500.00, c_продаж,     m_переказ,   'Продав старий ноутбук',         NOW() - INTERVAL '8 days'),
    (wid, false,  3400.00, c_електроніка,m_кредитна,  'Механічна клавіатура',          NOW() - INTERVAL '9 days'),
    (wid, false,   490.00, c_спорт,      m_дебетова,  'Абонемент у спортзал',          NOW() - INTERVAL '10 days'),
    (wid, false,   650.00, c_медицина,   m_готівка,   'Прийом у стоматолога',          NOW() - INTERVAL '11 days'),
    (wid, true,    420.00, c_відсотки,   m_переказ,   'Відсотки по депозиту',          NOW() - INTERVAL '12 days'),
    (wid, false,   299.00, c_телефон,    m_мобільний, 'Мобільний тариф',               NOW() - INTERVAL '14 days'),

  -- ── Минулий місяць ────────────────────────────────────────────
    (wid, true,  48000.00, c_зарплата,   m_переказ,   'Зарплата за місяць',            NOW() - INTERVAL '35 days'),
    (wid, false,  2200.00, c_продукти,   m_дебетова,  'АТБ + Сільпо',                  NOW() - INTERVAL '36 days'),
    (wid, false,   850.00, c_ресторан,   m_кредитна,  'Вечеря з другом',               NOW() - INTERVAL '38 days'),
    (wid, false,  1400.00, c_комунальні, m_переказ,   'Квартплата + електрика',        NOW() - INTERVAL '40 days'),
    (wid, false,   680.00, c_таксі,      m_мобільний, 'Таксі на вихідних',             NOW() - INTERVAL '42 days'),
    (wid, false,  1200.00, c_освіта,     m_кредитна,  'Онлайн-курс з дизайну',        NOW() - INTERVAL '44 days'),
    (wid, false,   450.00, c_спорт,      m_дебетова,  'Абонемент у спортзал',          NOW() - INTERVAL '46 days'),
    (wid, false,   199.00, c_ігри,       m_мобільний, 'Підписка Spotify + Netflix',    NOW() - INTERVAL '48 days'),
    (wid, false,   320.00, c_кафе,       m_готівка,   'Бранч у суботу',               NOW() - INTERVAL '50 days'),

  -- ── 2–3 місяці тому ───────────────────────────────────────────
    (wid, true,  48000.00, c_зарплата,   m_переказ,   'Зарплата за місяць',            NOW() - INTERVAL '65 days'),
    (wid, false,  1800.00, c_продукти,   m_дебетова,  'Новус + ринок',                 NOW() - INTERVAL '67 days'),
    (wid, false,  7200.00, c_медицина,   m_кредитна,  'Стоматолог — пломби',           NOW() - INTERVAL '69 days'),
    (wid, false,  2600.00, c_освіта,     m_переказ,   'Курс на Udemy (річна підписка)',NOW() - INTERVAL '71 days'),
    (wid, false,   890.00, c_ігри,       m_мобільний, 'Steam — декілька ігор',         NOW() - INTERVAL '73 days'),
    (wid, false,  1400.00, c_комунальні, m_переказ,   'Квартплата + електрика',        NOW() - INTERVAL '75 days'),
    (wid, false,   750.00, c_транспорт,  m_дебетова,  'Проїзний на місяць',            NOW() - INTERVAL '78 days'),
    (wid, true,  12000.00, c_продаж,     m_переказ,   'Продаж речей на OLX',           NOW() - INTERVAL '80 days'),

  -- ── 4–6 місяців тому ──────────────────────────────────────────
    (wid, true,  48000.00, c_зарплата,   m_переказ,   'Зарплата за місяць',            NOW() - INTERVAL '125 days'),
    (wid, false,  9500.00, c_електроніка,m_кредитна,  'Навушники Sony WH-1000XM5',    NOW() - INTERVAL '128 days'),
    (wid, false,  1400.00, c_комунальні, m_переказ,   'Квартплата + електрика',        NOW() - INTERVAL '132 days'),
    (wid, true,  48000.00, c_зарплата,   m_переказ,   'Зарплата за місяць',            NOW() - INTERVAL '160 days'),
    (wid, false,  6800.00, c_медицина,   m_готівка,   'Огляд у лікаря + аналізи',     NOW() - INTERVAL '163 days'),
    (wid, false,  3100.00, c_спорт,      m_дебетова,  'Новий велосипедний шолом',      NOW() - INTERVAL '167 days'),
    (wid, true,  48000.00, c_зарплата,   m_переказ,   'Зарплата за місяць',            NOW() - INTERVAL '190 days'),
    (wid, false,  1400.00, c_комунальні, m_переказ,   'Квартплата + електрика',        NOW() - INTERVAL '193 days'),
    (wid, true,    980.00, c_відсотки,   m_переказ,   'Відсотки по депозиту',          NOW() - INTERVAL '195 days');

END $$;
