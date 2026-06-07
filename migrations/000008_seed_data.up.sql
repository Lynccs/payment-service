-- ── Payment categories ──────────────────────────────────────────
INSERT INTO payment_categories (name, is_income) VALUES
  ('Продукти',               false),
  ('Кафе та ресторани',      false),
  ('Транспорт',              false),
  ('Комунальні послуги',     false),
  ('Одяг та взуття',         false),
  ('Здоровя та аптека',      false),
  ('Розваги',                false),
  ('Освіта',                 false),
  ('Техніка та електроніка', false),
  ('Спорт',                  false),
  ('Зарплата',               true),
  ('Фріланс',                true),
  ('Подарунок',              true),
  ('Кешбек',                 true);

-- ── Payment methods ──────────────────────────────────────────────
INSERT INTO payment_methods (name) VALUES
  ('Monobank'),
  ('ПриватБанк'),
  ('Готівка'),
  ('Apple Pay');

-- ── Payments ─────────────────────────────────────────────────────
-- Прив'язуємо до гаманця ОСТАННЬОГО зареєстрованого користувача.
-- Запускай після реєстрації.

DO $$
DECLARE
  wid       INTEGER;
  cat_id    INTEGER;
  mth_id    INTEGER;

  -- категорії витрат
  c_продукти    INTEGER;
  c_кафе        INTEGER;
  c_транспорт   INTEGER;
  c_комунальні  INTEGER;
  c_одяг        INTEGER;
  c_здоровя     INTEGER;
  c_розваги     INTEGER;
  c_освіта      INTEGER;
  c_техніка     INTEGER;
  c_спорт       INTEGER;
  -- категорії доходів
  c_зарплата    INTEGER;
  c_фріланс     INTEGER;
  c_подарунок   INTEGER;
  c_кешбек      INTEGER;
  -- методи
  m_mono        INTEGER;
  m_privat      INTEGER;
  m_cash        INTEGER;
  m_apple       INTEGER;
BEGIN
  SELECT id INTO wid FROM wallets ORDER BY id DESC LIMIT 1;
  IF wid IS NULL THEN
    RAISE EXCEPTION 'Спочатку зареєструйся — гаманець ще не створено';
  END IF;

  -- завантажуємо ID категорій
  SELECT id INTO c_продукти   FROM payment_categories WHERE name = 'Продукти'               AND is_income = false;
  SELECT id INTO c_кафе       FROM payment_categories WHERE name = 'Кафе та ресторани'      AND is_income = false;
  SELECT id INTO c_транспорт  FROM payment_categories WHERE name = 'Транспорт'              AND is_income = false;
  SELECT id INTO c_комунальні FROM payment_categories WHERE name = 'Комунальні послуги'     AND is_income = false;
  SELECT id INTO c_одяг       FROM payment_categories WHERE name = 'Одяг та взуття'         AND is_income = false;
  SELECT id INTO c_здоровя    FROM payment_categories WHERE name = 'Здоровя та аптека'      AND is_income = false;
  SELECT id INTO c_розваги    FROM payment_categories WHERE name = 'Розваги'                AND is_income = false;
  SELECT id INTO c_освіта     FROM payment_categories WHERE name = 'Освіта'                 AND is_income = false;
  SELECT id INTO c_техніка    FROM payment_categories WHERE name = 'Техніка та електроніка' AND is_income = false;
  SELECT id INTO c_спорт      FROM payment_categories WHERE name = 'Спорт'                  AND is_income = false;
  SELECT id INTO c_зарплата   FROM payment_categories WHERE name = 'Зарплата'               AND is_income = true;
  SELECT id INTO c_фріланс    FROM payment_categories WHERE name = 'Фріланс'                AND is_income = true;
  SELECT id INTO c_подарунок  FROM payment_categories WHERE name = 'Подарунок'              AND is_income = true;
  SELECT id INTO c_кешбек     FROM payment_categories WHERE name = 'Кешбек'                 AND is_income = true;

  -- завантажуємо ID методів
  SELECT id INTO m_mono   FROM payment_methods WHERE name = 'Monobank';
  SELECT id INTO m_privat FROM payment_methods WHERE name = 'ПриватБанк';
  SELECT id INTO m_cash   FROM payment_methods WHERE name = 'Готівка';
  SELECT id INTO m_apple  FROM payment_methods WHERE name = 'Apple Pay';

  -- перевірка що всі ID знайдено
  IF c_зарплата IS NULL OR m_mono IS NULL THEN
    RAISE EXCEPTION 'Категорії або методи не знайдено — перевір INSERT вище';
  END IF;

  -- ── Поточний місяць ───────────────────────────────────────────
  INSERT INTO payments (wallet_id, is_income, amount, category_id, payment_method_id, description, transaction_date) VALUES
    (wid, true,  42000.00, c_зарплата,   m_mono,   'Зарплата за місяць',          NOW() - INTERVAL '2 days'),
    (wid, false,  1850.00, c_продукти,   m_mono,   'Сільпо — тижневий запас',     NOW() - INTERVAL '1 day'),
    (wid, false,   420.00, c_кафе,       m_apple,  'Обід з колегами',             NOW() - INTERVAL '3 days'),
    (wid, false,   680.00, c_транспорт,  m_mono,   'Таксі + поповнення карти метро', NOW() - INTERVAL '4 days'),
    (wid, false,  1200.00, c_комунальні, m_privat, 'Квартплата + світло',         NOW() - INTERVAL '5 days'),
    (wid, true,   8500.00, c_фріланс,    m_mono,   'Проєкт для клієнта',          NOW() - INTERVAL '6 days'),
    (wid, false,   990.00, c_розваги,    m_apple,  'Кінотеатр + попкорн',         NOW() - INTERVAL '7 days'),
    (wid, false,  3200.00, c_одяг,       m_mono,   'Кросівки Nike',               NOW() - INTERVAL '9 days'),
    (wid, false,   560.00, c_здоровя,    m_cash,   'Ліки',                        NOW() - INTERVAL '10 days'),
    (wid, true,    350.00, c_кешбек,     m_mono,   'Кешбек Monobank',             NOW() - INTERVAL '11 days'),

  -- ── Минулий місяць ────────────────────────────────────────────
    (wid, true,  42000.00, c_зарплата,   m_mono,   'Зарплата за місяць',          NOW() - INTERVAL '35 days'),
    (wid, false,  2100.00, c_продукти,   m_mono,   'АТБ + Сільпо',                NOW() - INTERVAL '36 days'),
    (wid, false,   780.00, c_кафе,       m_apple,  'Суші-бар з другом',           NOW() - INTERVAL '38 days'),
    (wid, false,  4500.00, c_техніка,    m_privat, 'Навушники Sony',              NOW() - INTERVAL '40 days'),
    (wid, false,  1200.00, c_комунальні, m_privat, 'Квартплата + світло',         NOW() - INTERVAL '42 days'),
    (wid, false,   650.00, c_спорт,      m_mono,   'Абонемент у спортзал',        NOW() - INTERVAL '44 days'),
    (wid, true,   5000.00, c_подарунок,  m_cash,   'День народження',             NOW() - INTERVAL '48 days'),
    (wid, false,   310.00, c_транспорт,  m_cash,   'Маршрутка + таксі',           NOW() - INTERVAL '50 days'),

  -- ── 2–3 місяці тому ───────────────────────────────────────────
    (wid, true,  42000.00, c_зарплата,   m_mono,   'Зарплата за місяць',          NOW() - INTERVAL '65 days'),
    (wid, false,  1750.00, c_продукти,   m_privat, 'Новус + ринок',               NOW() - INTERVAL '67 days'),
    (wid, false,  2400.00, c_освіта,     m_mono,   'Курс на Udemy',               NOW() - INTERVAL '70 days'),
    (wid, true,  12000.00, c_фріланс,    m_mono,   'Великий проєкт',              NOW() - INTERVAL '72 days'),
    (wid, false,   890.00, c_розваги,    m_apple,  'Steam ігри',                  NOW() - INTERVAL '75 days'),
    (wid, false,  1200.00, c_комунальні, m_privat, 'Квартплата + світло',         NOW() - INTERVAL '78 days'),

  -- ── 4–6 місяців тому ──────────────────────────────────────────
    (wid, true,  42000.00, c_зарплата,   m_mono,   'Зарплата за місяць',          NOW() - INTERVAL '125 days'),
    (wid, false,  8900.00, c_одяг,       m_mono,   'Зимова куртка + черевики',    NOW() - INTERVAL '130 days'),
    (wid, true,  42000.00, c_зарплата,   m_mono,   'Зарплата за місяць',          NOW() - INTERVAL '160 days'),
    (wid, false,  6500.00, c_здоровя,    m_privat, 'Стоматолог',                  NOW() - INTERVAL '165 days'),
    (wid, true,  42000.00, c_зарплата,   m_mono,   'Зарплата за місяць',          NOW() - INTERVAL '190 days'),
    (wid, false,  3800.00, c_техніка,    m_mono,   'Механічна клавіатура',        NOW() - INTERVAL '195 days');

END $$;
