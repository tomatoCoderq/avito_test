import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString, randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// Конфигурация теста
export const options = {
  // Пороговые значения для проверки производительности
  thresholds: {
    // 95% запросов должны быть быстрее 100мс
    http_req_duration: ['p(95) < 100'],
    // Не должно быть ошибок HTTP
    http_req_failed: ['rate < 0.01'],
  },
  
  // Сценарий нагрузочного тестирования (общая длительность 30 секунд)
  scenarios: {
    deactivate_users_stress: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages: [
        { duration: '5s', target: 10 },   // Разгон до 10 пользователей
        { duration: '20s', target: 20 },  // Поддержание 20 пользователей
        { duration: '5s', target: 0 },    // Плавное снижение
      ],
    },
  },
};

// Глобальные переменные для хранения созданных данных
let createdTeams = [];

// Модуль setup - подготовка тестовых данных
export function setup() {
  console.log('🚀 Начинаем подготовку тестовых данных...');
  
  const baseUrl = 'http://localhost:8080';
  const teams = [];
  const teamsToCreate = 50;
  const usersPerTeam = 30;
  
  // Создаем 50 команд с 30 уникальными пользователями в каждой
  for (let teamIndex = 0; teamIndex < teamsToCreate; teamIndex++) {
    const teamName = `stress_team_${teamIndex}_${randomString(6)}`;
    const members = [];
    
    // Генерируем 30 уникальных пользователей для команды
    for (let userIndex = 0; userIndex < usersPerTeam; userIndex++) {
      members.push({
        user_id: `stress_user_${teamIndex}_${userIndex}_${randomString(4)}`,
        username: `StressUser${teamIndex}_${userIndex}`,
        is_active: true
      });
    }
    
    const teamData = {
      team_name: teamName,
      members: members
    };
    
    // Отправляем запрос на создание команды
    const response = http.post(`${baseUrl}/team/add`, JSON.stringify(teamData), {
      headers: { 'Content-Type': 'application/json' },
    });
    
    if (check(response, {
      'team created successfully': (r) => r.status === 200 || r.status === 201,
    })) {
      teams.push({
        name: teamName,
        userIds: members.map(member => member.user_id)
      });
      
      if ((teamIndex + 1) % 10 === 0) {
        console.log(`✅ Создано ${teamIndex + 1}/${teamsToCreate} команд`);
      }
    } else {
      console.error(`❌ Ошибка создания команды ${teamName}: ${response.status}`);
    }
    
    // Небольшая пауза между созданием команд, чтобы не перегружать сервер
    sleep(0.1);
  }
  
  console.log(`🎉 Подготовка завершена! Создано ${teams.length} команд с ${teams.length * usersPerTeam} пользователями`);
  
  return {
    baseUrl: baseUrl,
    teams: teams
  };
}

// Основная функция нагрузочного тестирования
export default function (data) {
  if (!data || !data.teams || data.teams.length === 0) {
    console.error('❌ Нет подготовленных данных для тестирования');
    return;
  }
  
  // Выбираем случайную команду
  const randomTeam = randomItem(data.teams);
  
  // Выбираем случайное количество пользователей для деактивации (от 1 до 10)
  const usersToDeactivate = Math.floor(Math.random() * 10) + 1;
  const selectedUserIds = randomTeam.userIds
    .sort(() => 0.5 - Math.random())
    .slice(0, usersToDeactivate);
  
  const deactivationData = {
    team_name: randomTeam.name,
    user_ids: selectedUserIds
  };
  
  // Выполняем запрос на деактивацию пользователей
  const response = http.post(
    `${data.baseUrl}/team/deactivateUsers`,
    JSON.stringify(deactivationData),
    {
      headers: { 'Content-Type': 'application/json' },
      tags: { 
        scenario: 'deactivate_users',
        users_count: selectedUserIds.length.toString()
      },
    }
  );
  
  // Проверяем результат запроса
  check(response, {
    'deactivation request successful': (r) => r.status === 200 || r.status === 201,
    'response time under 100ms': (r) => r.timings.duration < 100,
    'response has body': (r) => r.body && r.body.length > 0,
  });
  
  // Логируем медленные запросы
  if (response.timings.duration > 100) {
    console.warn(`⚠️ Медленный запрос: ${response.timings.duration}ms для ${selectedUserIds.length} пользователей`);
  }
  
  // Короткая пауза между запросами
  sleep(0.1);
}

// Очистка после тестирования
export function teardown(data) {
  console.log('🧹 Очистка тестовых данных завершена');
  
  if (data && data.teams) {
    console.log(`📊 Было протестировано ${data.teams.length} команд`);
  }
}

